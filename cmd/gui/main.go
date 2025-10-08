package main

import (
	"database/sql"
	t "vado/internal/domain/task"
	u "vado/internal/domain/user"
	"vado/internal/gui"
	"vado/internal/server"
	"vado/internal/server/context"
	"vado/pkg/logger"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()

	db := server.InitDB()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Инициализируем сервисы
	userService := u.NewUserService(u.NewUserDBRepo(db))
	taskService := t.NewTaskService(t.NewTaskDBRepo(db))

	appCtx := &context.AppContext{
		DB:     db,
		Logger: log,
		HTTP:   context.CreateHTTPContext(userService, taskService),
		GRPC:   server.NewServerGRPC(userService, taskService, ":50051"),
	}

	gui.ShowApp(appCtx)
}
