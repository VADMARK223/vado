package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"vado/internal/util"
	"vado/pkg/logger"

	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("Hello, world!")
	log1, _ := logger.Init()
	defer logger.Sync()

	log1.Info(fmt.Sprintf("Starting CLI mode. (%s)", util.GetModeValue()))

	// Настройки подключения
	broker := "localhost:9092"
	topic := "tasks"

	// ======== 1. Producer ========
	brokers := []string{broker}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "tasks",
		Balancer:               &kafka.LeastBytes{}, // Балансировщик для распределения сообщений по партициям (можно использовать другие: Hash, RoundRobin)
		AllowAutoTopicCreation: true,                // Авто создание топика
	}
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			println(err)
		}
	}(writer)

	msg := kafka.Message{
		Key:   []byte("task-4"),
		Value: []byte(`{"id":1,"name":"Do homework","done":false}`),
		Time:  time.Now(),
	}

	fmt.Println("🚀 Отправляем сообщение в Kafka...")
	if err := writer.WriteMessages(context.Background(), msg); err != nil {
		log.Fatalf("❌ Ошибка отправки: %v", err)
	}
	fmt.Println("✅ Сообщение отправлено")

	// ======== 2. Consumer ========
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "test-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			println(err)
		}
	}(reader)

	fmt.Println("👂 Читаем сообщение из Kafka...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Fatalf("❌ Ошибка чтения: %v", err)
	}

	fmt.Printf("📩 Получено сообщение:\n  key=%s\n  value=%s\n", string(m.Key), string(m.Value))

	//startServer()
}

//func startServer() {
//	var r repo.TaskRepo
//	r = repo.NewTaskDBRepo(constant.GetDSN())
//	var s service.ITaskService = service.NewTaskService(r)
//	err := http.StartHTTPServer(s)
//	if err != nil {
//		fmt.Println(fmt.Sprintf("Error start server: %s", err.Error()))
//	}
//}
