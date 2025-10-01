package kafka

import (
	"fmt"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func CheckKafkaConnection() error {
	conn, err := kafka.Dial("tcp", "localhost:9092")
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
