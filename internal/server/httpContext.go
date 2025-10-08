package server

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"vado/internal/domain/user"
	"vado/internal/service/task"
	http2 "vado/internal/transport/rest"
	"vado/internal/util"
	"vado/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

var httpMtx sync.Mutex

func InitHTTPContext(userService *user.Service, taskService *task.Service) (*util.HttpContext, error) {
	httpMtx.Lock()
	httpCtx := &util.HttpContext{}
	httpCtx.UserService = userService
	httpCtx.TaskService = taskService

	if httpCtx.ServerHTTP != nil {
		return httpCtx, errors.New("serverHTTP already running")
	}

	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("data/index.html"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := PageData{
			Title:   "Привет из Vado 🚀",
			Message: fmt.Sprintf("Сервер работает. (%s)", strings.ToUpper(util.GetModeValue())),
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			return
		}
	})

	taskHandler := &http2.TaskHandler{Service: taskService}
	mux.HandleFunc("/tasks", taskHandler.TasksHandler)
	mux.HandleFunc("/tasks/", taskHandler.TaskByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/slow", taskHandler.SlowHandler)

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		newUser := user.User{Username: "test", Password: "test"}
		err := userService.CreateUser(newUser)
		if err != nil {
			logger.L().Error("create user failed", zap.Error(err))
			return
		}
		_, _ = w.Write([]byte("USER CREATED!"))
	})

	httpCtx.ServerHTTP = &http.Server{
		Addr:    ":5556",
		Handler: mux,
	}
	httpMtx.Unlock()

	logger.L().Info("HTTP-serverHTTP started on :5556")

	// ListenAndServe блокирующий
	// ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
	// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: rest: Server closed» даже при нормальной остановке.
	go func() {
		if err := httpCtx.ServerHTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Error("HTTP serverHTTP error", zap.Error(err))
		}
	}()

	return httpCtx, nil
}

type PageData struct {
	Title   string
	Message string
}
