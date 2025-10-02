package kafka

import (
	"context"
	"log"
	"time"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Produce(message string) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "tasks",
		Balancer: &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
	}
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			logger.L().Error("Error close writer:", zap.Error(err))
		}
	}(writer)

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339)),
			Value: []byte(message),
		},
	)
	if err != nil {
		logger.L().Error("Failed to write message:", zap.Error(err))
		return
	}
	log.Println("Message sent to Kafka:", message)

	logger.L().Info("Message sent to Kafka:", zap.String("message", message))
}
