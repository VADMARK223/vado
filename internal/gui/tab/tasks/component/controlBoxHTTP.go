package component

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"image/color"
	"net/http"
	"sync"
	"time"
	"vado/internal/gui/common"
	constant2 "vado/internal/gui/constant"
	"vado/internal/service"
	http2 "vado/internal/transport/rest"
	"vado/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	guiUpdateMillisecond = 500 // Частота обновления GUI
	shutdownSecond       = 5
)

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	httpMtx       sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func NewControlBoxHTTP(service service.ITaskService) fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер HTTP:")
	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), nil)
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

	go func() {
		for {
			httpMtx.Lock()
			running := srv == nil
			inProcess := stopInProcess
			httpMtx.Unlock()
			if running {
				fyne.Do(func() {
					waitLbl.Hide()
					startBtn.Enable()
					stopBtn.Disable()
					statusIndicator.SetFillColor(constant2.Red())
				})
			} else {
				fyne.Do(func() {
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
				})
			}

			time.Sleep(time.Millisecond * guiUpdateMillisecond)
		}
	}()

	if util.AutoStartServer() {
		startOnTapped(service)
	}

	return container.NewHBox(lbl, startBtn, stopBtn, container.NewCenter(statusIndicator), waitLbl)
}

func startOnTapped(service service.ITaskService) {
	go func() {
		if err := StartServer(service); err != nil {
			fmt.Println("Start server error:", err)
		}
	}()
}

type PageData struct {
	Title   string
	Message string
}

func StartServer(service service.ITaskService) error {
	httpMtx.Lock()
	if srv != nil {
		return errors.New("server already running")
	}

	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	handler := &http2.TaskHandler{Service: service}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//_, _ = w.Write([]byte("Hello from vado!"))
		tmpl := template.Must(template.ParseFiles("data/index.html"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := PageData{
			Title:   "Привет из Vado 🚀",
			Message: "Сервер работает.",
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			return
		}
	})
	mux.HandleFunc("/tasks", handler.TasksHandler)
	mux.HandleFunc("/tasks/", handler.TaskByIDHandler)
	mux.HandleFunc("/slow", handler.SlowHandler)

	srv = &http.Server{
		Addr:    ":5555",
		Handler: mux,
	}
	httpMtx.Unlock()

	fmt.Println("HTTP-server started on :5555")

	// ListenAndServe блокирующий
	// ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
	// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: rest: Server closed» даже при нормальной остановке.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error:", err)
	}

	return nil
}

func stopServer() {
	httpMtx.Lock()
	s := srv
	if s == nil {
		panic("Server is nil!")
	}
	httpMtx.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownSecond*time.Second)
	defer func() {
		fmt.Println("Context canceled.")
		cancel()
	}()

	// Блокирующая операция
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	} else {
		fmt.Println("Server stopped.")
	}
	httpMtx.Lock()
	stopInProcess = false
	srv = nil
	httpMtx.Unlock()
}
