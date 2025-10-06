package http

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"image/color"
	"net/http"
	"strings"
	"sync"
	"time"
	"vado/internal/gui/common"
	constant2 "vado/internal/gui/constant"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/service"
	http2 "vado/internal/transport/rest"
	"vado/internal/util"
	"vado/pkg/logger"

	_ "vado/internal/docs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const (
	shutdownSecond = 5
)

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	httpMtx       sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func NewControlBoxHTTP(service service.ITaskService) fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер HTTP:")
	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), nil)
	startBtn.Disable()
	startBtn.OnTapped = func() {
		startOnTapped(service)
	}

	stopBtn := common.NewBtn("Стоп", theme.MediaStopIcon(), func() {
		httpMtx.Lock()
		if srv == nil || stopInProcess {
			httpMtx.Unlock()
			return
		}
		stopInProcess = true
		httpMtx.Unlock()
		go stopServer()
	})

	statusIndicator := common.NewIndicator(color.RGBA{R: 255, G: 0, B: 0, A: 255}, fyne.NewSize(15, 15))

	waitLbl := widget.NewLabel("Остановка сервера...")
	waitLbl.Hide()

	//ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(time.Millisecond * component.GuiUpdateMillisecond)
		for {
			select {
			//case <-ctx.Done():
			//	return
			case <-ticker.C:
				httpMtx.Lock()
				running := srv == nil
				inProcess := stopInProcess
				httpMtx.Unlock()

				fyne.Do(func() {
					if running {
						waitLbl.Hide()
						startBtn.Enable()
						stopBtn.Disable()
						statusIndicator.SetFillColor(constant2.Red())
					} else {
						if inProcess {
							waitLbl.Show()
							startBtn.Disable()
							stopBtn.Disable()
							statusIndicator.SetFillColor(constant2.Orange())
						} else {
							waitLbl.Hide()
							startBtn.Disable()
							stopBtn.Enable()
							statusIndicator.SetFillColor(constant2.Green())
						}
					}

				})
			}
		}
	}()

	if util.AutoStartServerHTTP() {
		startOnTapped(service)
	}

	//time.Sleep(time.Second * 5)
	//cancel()
	return container.NewHBox(lbl, startBtn, stopBtn, container.NewCenter(statusIndicator), waitLbl)
}

func startOnTapped(service service.ITaskService) {
	go func() {
		if err := StartHTTPServer(service); err != nil {
			fmt.Println("Start server error:", err)
		}
	}()
}

type PageData struct {
	Title   string
	Message string
}

func StartHTTPServer(service service.ITaskService) error {
	httpMtx.Lock()
	if srv != nil {
		return errors.New("server already running")
	}

	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	handler := &http2.TaskHandler{Service: service}
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
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/", handler.TaskByIDHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/slow", handler.SlowHandler)

	srv = &http.Server{
		Addr:    ":5556",
		Handler: mux,
	}
	httpMtx.Unlock()

	logger.L().Info("HTTP-server started on :5556")

	// ListenAndServe блокирующий
	// ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
	// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: rest: Server closed» даже при нормальной остановке.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.L().Error("HTTP server error", zap.Error(err))
	}

	return nil
}

func stopServer() {
	httpMtx.Lock()
	s := srv
	httpMtx.Unlock()

	if s == nil {
		logger.L().Warn("Server already stopped")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownSecond*time.Second)
	defer func() {
		logger.L().Debug("Context canceled.")
		cancel()
	}()

	// Блокирующая операция
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	} else {
		logger.L().Info("HTTP server stopped.")
	}
	httpMtx.Lock()
	stopInProcess = false
	srv = nil
	httpMtx.Unlock()
}
