package main

import (
	"database/sql"
	t "vado/internal/domain/task"
	u "vado/internal/domain/user"
	"vado/internal/gui"
	"vado/internal/server"
	"vado/internal/server/context"
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

	// Инициализируем сервисы
	userService := u.NewUserService(u.NewUserDBRepo(db))
	taskService := t.NewTaskService(t.NewTaskDBRepo(db))

	// Запускаем gRPC
	grpcServer := server.NewServerGRPC(userService, taskService, ":50051")
	if err := grpcServer.Start(); err != nil {
		logger.L().Error("gRPC server error", zap.Error(err))
	}
	defer grpcServer.Stop()

	httpCtx := context.CreateHTTPContext(userService, taskService)
	err := httpCtx.Start()
	if err != nil {
		logger.L().Error("HTTP server error", zap.Error(err))
		return
	}
	if httpCtx.ServerHTTP != nil {
		defer func() {
			_ = httpCtx.ServerHTTP.Close()
		}()
	}

	appCtx := &context.AppContext{
		DB:          db,
		Logger:      log,
		HttpContext: httpCtx,
		GRPC:        grpcServer,
	}

	gui.ShowMainApp(appCtx)
}
