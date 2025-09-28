package main

import (
	"fmt"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/repo"
	"vado/internal/repo/db"
	"vado/internal/service"
)

func main() {
	startServer()
}

func startServer() {
	fmt.Println("Vado start...")
	var r repo.TaskRepo
	r = db.NewTaskDBRepo(constant.TasksDataSourceName)
	var s service.ITaskService = service.NewTaskService(r)
	err := component.StartServer(s)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error start server: %s", err.Error()))
	}
}
