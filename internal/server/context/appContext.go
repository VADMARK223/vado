package context

import (
	"database/sql"
	"net/http"
	"sync"
	"vado/internal/domain/task"
	"vado/internal/domain/user"
	"vado/internal/server"

	"go.uber.org/zap"
)

type AppContext struct {
	DB     *sql.DB
	Logger *zap.Logger
	HTTP   *HttpContext
	GRPC   *server.GRPCServer
}

type HttpContext struct {
	ServerHTTP  *http.Server
	UserService *user.Service
	TaskService *task.Service

	running bool
	mu      sync.RWMutex
}
