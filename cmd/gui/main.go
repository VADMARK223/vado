package main

import (
	"database/sql"
	task2 "vado/internal/domain/task"
	user2 "vado/internal/domain/user"
	"vado/internal/gui"
	"vado/internal/server"
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

	userService := user2.NewUserService(user2.NewUserDBRepo(db))
	taskService := task2.NewTaskService(task2.NewTaskDBRepo(db))
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
