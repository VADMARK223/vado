package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *Config) *Producer {
	/*return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  cfg.Brokers,
			Topic:    cfg.Topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}*/

	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(cfg.Brokers...),
			Topic:                  "tasks",
			Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
			AllowAutoTopicCreation: true,                // Авто создание топика
		},
	}
}

func (p *Producer) Close() {
	_ = p.writer.Close()
}

func (p *Producer) Send(key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}
	if err := p.writer.WriteMessages(context.Background(), msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	log.Printf("✅ sent message: key=%s value=%s", key, value)
	return nil
}
