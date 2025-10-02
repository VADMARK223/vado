package kafka

import (
	"fmt"
	"net"
	"os"
	"strings"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// GetKafkaBrokers возвращает список брокеров из переменной окружения или localhost
func GetKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// Автоопределение: если можем резолвить 'kafka', значит мы в Docker
		brokers = detectKafkaAddress()

	}
	return strings.Split(brokers, ",")
}

// GetKafkaBroker возвращает первый брокер (для Dial)
func GetKafkaBroker() string {
	return GetKafkaBrokers()[0]
}

func detectKafkaAddress() string {
	// Пробуем резолвить имя 'kafka' (работает только в Docker)
	_, err := net.LookupHost("kafka")
	if err == nil {
		// Мы в Docker-сети, используем Docker DNS
		logger.L().Info("Detected Docker environment, using 'kafka:9092'")
		return "kafka:9092"
	}

	// Мы на хосте, используем localhost
	logger.L().Info("Detected host environment, using 'localhost:9092'")
	return "localhost:9092"
}

func CheckKafkaConnection() error {
	broker := GetKafkaBroker()
	logger.L().Info("Connecting to Kafka", zap.String("broker", broker))

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer func(conn *kafka.Conn) {
		_ = conn.Close()
	}(conn)

	// Получаем список брокеров
	brokers, err := conn.Brokers()
	if err != nil {
		return fmt.Errorf("failed to get brokers: %w", err)
	}

	logger.L().Info("Kafka connection successful",
		zap.Int("brokers_count", len(brokers)),
		zap.Any("brokers", brokers))
	return nil
}
