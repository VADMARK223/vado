package server

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"vado/internal/model"
	"vado/internal/service/user"
	"vado/internal/util"
	"vado/pkg/logger"

	"go.uber.org/zap"
)

var (
	//serverHTTP *http.Server
	httpMtx sync.Mutex
)

// func InitHTTPContext(taskService service.ITaskService, userService *user.UserService) (*util.HttpContext, error) {
func InitHTTPContext(userService *user.Service) (*util.HttpContext, error) {
	httpMtx.Lock()
	httpCtx := &util.HttpContext{}
	httpCtx.UserService = userService

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

	/*handler := &http2.TaskHandler{Service: taskService}
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/", handler.TaskByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/slow", handler.SlowHandler)

	*/

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		newUser := model.User{Username: "test", Password: "test"}
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
