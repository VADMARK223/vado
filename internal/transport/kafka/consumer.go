package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *Config) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			Topic:    cfg.Topic,
			GroupID:  cfg.GroupID,
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) Close() {
	_ = c.reader.Close()
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Printf("👂 consuming from topic %q ...", c.reader.Config().Topic)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		log.Printf("📩 received: key=%s value=%s", string(m.Key), string(m.Value))
	}
}

/*import (
	"context"
	"log"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func main() {
	brokers := []string{"localhost:9092"}
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
}*/
