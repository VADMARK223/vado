package main

import (
	"fmt"
	"time"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/repo"
	"vado/internal/repo/db"
	"vado/internal/service"
	"vado/internal/transport/kafka"
	"vado/internal/util"
	"vado/pkg/logger"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()

	log.Info(fmt.Sprintf("Starting CLI mode. (%s)", util.GetModeValue()))

	currentTime := time.Now().String()[:19]
	go kafka.Consume() // Запускаем в фоне consumer

	message := fmt.Sprintf("Message: %s", currentTime)
	kafka.Produce(message)

	startServer()
}

func startServer() {
	var r repo.TaskRepo
	r = db.NewTaskDBRepo(constant.GetDSN())
	var s service.ITaskService = service.NewTaskService(r)
	err := component.StartServer(s)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error start server: %s", err.Error()))
	}
}
