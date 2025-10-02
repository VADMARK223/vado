package kafka

import (
	"context"
	"time"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Produce(message string) {
	brokers := GetKafkaBrokers()
	logger.L().Info("Kafka producer starting", zap.Strings("brokers", brokers))

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "tasks",
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Автосоздание топика
	}
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			logger.L().Error("Error close writer:", zap.Error(err))
		}
	}(writer)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339)),
			Value: []byte(message),
		},
	)
	if err != nil {
		logger.L().Error("Failed to write message:", zap.Error(err))
		return
	}
	logger.L().Info("Message sent to Kafka:", zap.String("message", message))
}
