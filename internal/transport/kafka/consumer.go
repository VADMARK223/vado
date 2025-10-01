package kafka

import (
	"context"
	"log"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Consume() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "tasks",
		GroupID: "task-consumers",
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			logger.L().Error("Error close reader:", zap.Error(err))
		}
	}(reader)

	logger.L().Info("Kafka consumer started...")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.L().Error("Error reading message:", zap.Error(err))
			break
		}
		logger.L().Info("Kafka message received:", zap.ByteString("message", m.Value), zap.String("key", string(m.Key)))
		log.Printf("Message received: key=%s value=%s", string(m.Key), string(m.Value))
	}
}
