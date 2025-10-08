package kafka

import "os"

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func Load() *Config {
	broker := getenv("KAFKA_BROKER", "localhost:9092")
	topic := getenv("KAFKA_TOPIC", "tasks")
	group := getenv("KAFKA_GROUP_ID", "go-consumer-group")

	return &Config{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: group,
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
