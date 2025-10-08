package main

import (
	"database/sql"
	"vado/internal/gui"
	task2 "vado/internal/repo/task"
	user2 "vado/internal/repo/user"
	"vado/internal/server"
	"vado/internal/service/task"
	"vado/internal/service/user"
	"vado/internal/util"
	"vado/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()

	db := server.InitDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	userService := user.NewUserService(user2.NewUserDBRepo(db))
	taskService := task.NewTaskService(task2.NewTaskDBRepo(db))
	http, err := server.InitHTTPContext(userService, taskService)
	if err != nil {
		logger.L().Error("Error init http server:", zap.Error(err))
		return
	}
	defer func() {
		_ = http.ServerHTTP.Close()
	}()

	appCtx := &util.AppContext{
		DB:          db,
		Logger:      log,
		HttpContext: http,
	}

	gui.ShowMainApp(appCtx)
}
