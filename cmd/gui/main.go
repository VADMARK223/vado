package main

import (
	"database/sql"
	"vado/internal/gui"
	user2 "vado/internal/repo/user"
	"vado/internal/server"
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

	testService := user.NewUserService(user2.NewUserDBRepo(db))
	http, err := server.InitHTTPContext(testService)
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

	/*taskService := service.NewTaskService(repo.NewTaskDBRepo(appCtx.DB))
	userService := service.NewUserService(repo.NewUserDBRepo(appCtx.DB))


	mux := http3.NewServeMux() // multiplexer = «распределитель запросов»
	mux.HandleFunc("/", func(w http3.ResponseWriter, r *http3.Request) {
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

	handler := &http2.TaskHandler{Service: taskService}
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/", handler.TaskByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/slow", handler.SlowHandler)

	mux.HandleFunc("/users", func(w http3.ResponseWriter, r *http3.Request) {
		user := model.User{Username: "admin1", Password: "ячс1"}
		err := userService.CreateUser(user)
		if err != nil {
			logger.L().Error("create user failed", zap.Error(err))
			return
		}
		_, _ = w.Write([]byte("USER CREATED!"))
	})
	http.ServerHTTP.Handler = mux*/

	gui.ShowMainApp(appCtx)
}

type PageData struct {
	Title   string
	Message string
}
