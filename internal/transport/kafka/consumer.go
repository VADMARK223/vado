package kafka

import (
	"context"
	"log"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Consume() {
	brokers := GetKafkaBrokers()
	logger.L().Info("Kafka consumer starting", zap.Strings("brokers", brokers))

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "tasks",
		GroupID: "task-consumers",
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			logger.L().Error("Error close reader:", zap.Error(err))
		}
	}(reader)

	logger.L().Info("Kafka consumer started and waiting for messages...")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.L().Error("Error reading message:", zap.Error(err))
			break
		}
		logger.L().Info("Kafka message received:",
			zap.ByteString("message", m.Value),
			zap.String("key", string(m.Key)),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset))
		log.Printf("Message received: key=%s value=%s", string(m.Key), string(m.Value))
	}
}
