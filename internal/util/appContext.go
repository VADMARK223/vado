package util

import (
	"database/sql"
	"net/http"
	"vado/internal/domain/task"
	"vado/internal/domain/user"

	"go.uber.org/zap"
)

type AppContext struct {
	DB          *sql.DB
	Logger      *zap.Logger
	HttpContext *HttpContext
}

type HttpContext struct {
	ServerHTTP  *http.Server
	UserService *user.Service
	TaskService *task.Service
}
