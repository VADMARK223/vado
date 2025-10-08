package context

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"vado/internal/domain/task"
	http2 "vado/internal/domain/task/transport/rest"
	"vado/internal/domain/user"
	"vado/internal/util"
	"vado/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

var httpMtx sync.Mutex

func CreateHTTPContext(userService *user.Service, taskService *task.Service) *HttpContext {
	httpMtx.Lock()
	defer httpMtx.Unlock()

	httpCtx := &HttpContext{
		UserService: userService,
		TaskService: taskService,
	}

	/*if err := httpCtx.Start(); err != nil {
		return nil, fmt.Errorf("failed to start HTTP server: %w", err)
	}*/

	return httpCtx
}

func (h *HttpContext) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return errors.New("HTTP server already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("data/index.html"))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := PageData{
			Title:   "Привет из Vado 🚀",
			Message: fmt.Sprintf("Сервер работает. (%s)", strings.ToUpper(util.GetModeValue())),
		}
		_ = tmpl.Execute(w, data)
	})

	// Остальные обработчики
	taskHandler := &http2.TaskHandler{Service: h.TaskService}
	mux.HandleFunc("/tasks", taskHandler.TasksHandler)
	mux.HandleFunc("/tasks/", taskHandler.TaskByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/slow", taskHandler.SlowHandler)
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		newUser := user.User{Username: "test", Password: "test"}
		if err := h.UserService.CreateUser(newUser); err != nil {
			logger.L().Error("create user failed", zap.Error(err))
			return
		}
		_, _ = w.Write([]byte("USER CREATED!"))
	})

	h.ServerHTTP = &http.Server{
		Addr:    ":5556",
		Handler: mux,
	}

	h.running = true
	logger.L().Info("HTTP server started on :5556")

	go func() {
		if err := h.ServerHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Error("HTTP server error", zap.Error(err))
		}
		h.mu.Lock()
		h.running = false
		h.mu.Unlock()
	}()

	return nil
}

func (h *HttpContext) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return
	}

	if err := h.ServerHTTP.Close(); err != nil {
		logger.L().Error("failed to stop HTTP server", zap.Error(err))
	}
	h.running = false
	logger.L().Info("HTTP server stopped")
}

func (h *HttpContext) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.running
}

type PageData struct {
	Title   string
	Message string
}
