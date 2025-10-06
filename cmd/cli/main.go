package main

import (
	"fmt"
	"time"
	"vado/internal/gui/tab/tasks/component/http"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/repo"
	"vado/internal/service"
	"vado/internal/transport/kafka"
	"vado/internal/util"
	"vado/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()

	log.Info(fmt.Sprintf("Starting CLI mode. (%s)", util.GetModeValue()))

	if err := kafka.CheckKafkaConnection(); err != nil {
		logger.L().Error("Kafka connection failed", zap.Error(err))
		return
	}

	currentTime := time.Now().String()[:19]
	go kafka.Consume() // Запускаем в фоне consumer
	message := fmt.Sprintf("Message: %s", currentTime)
	kafka.Produce(message)

	startServer()
}

func startServer() {
	var r repo.TaskRepo
	r = repo.NewTaskDBRepo(constant.GetDSN())
	var s service.ITaskService = service.NewTaskService(r)
	err := http.StartHTTPServer(s)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error start server: %s", err.Error()))
	}
}
