package main

import (
	"fmt"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/repo"
	"vado/internal/repo/db"
	"vado/internal/service"
	"vado/internal/util"
	"vado/pkg/logger"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()
	log.Info(fmt.Sprintf("Starting CLI mode. (%s)", util.GetModeValue()))
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
