package http

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"sync"
	"time"
	"vado/internal/gui/common"
	constant2 "vado/internal/gui/constant"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const autoStart = true            // Автоматически стартовать сервер
const guiUpdateMillisecond = 500  // Частота обновления GUI
const slowRequestDelaySecond = 10 // Длительность выполнения медленного запроса
const shutdownSecond = 5          // Время "мягкой" остановки сервера

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	httpMtx       sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func CreateView() fyne.CanvasObject {
	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), nil)
	startBtn.Disable()

	stopBtn := common.NewBtn("Стоп", theme.MediaStopIcon(), stopServer)
	stopBtn.Disable()

	waitLbl := widget.NewLabel("Остановка сервера...")
	waitLbl.Hide()

	statusLbl := widget.NewLabel("Состояние сервера:")

	statusIndicator := canvas.NewCircle(color.White)
	statusIndicator.FillColor = constant2.Red()
	statusIndicator.StrokeColor = color.Gray{Y: 0x99}
	statusIndicator.StrokeWidth = 1
	statusIndicatorLayout := container.NewWithoutLayout(statusIndicator)
	statusIndicator.Resize(fyne.NewSize(30, 30))

	startBtn.OnTapped = func() {
		go StartServer()
	}

	stopBtn.OnTapped = func() {
		httpMtx.Lock()
		if srv == nil || stopInProcess {
			httpMtx.Unlock()
			return
		}
		stopInProcess = true
		httpMtx.Unlock()
		go stopServer()
	}

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
					statusIndicator.FillColor = constant2.Red()
				})
			} else {
				fyne.Do(func() {
					if inProcess {
						waitLbl.Show()
						startBtn.Disable()
						stopBtn.Disable()
						statusIndicator.FillColor = constant2.Orange()
					} else {
						waitLbl.Hide()
						startBtn.Disable()
						stopBtn.Enable()
						statusIndicator.FillColor = constant2.Green()
					}
				})
			}

			time.Sleep(time.Millisecond * guiUpdateMillisecond)
		}
	}()

	controlBox := container.NewHBox(startBtn, stopBtn, waitLbl)
	mainVerticalBox := container.NewVBox(container.NewHBox(statusLbl, statusIndicatorLayout), controlBox)
	mainVerticalBox.Add(widget.NewSeparator())
	mainVerticalBox.Add(createMoneyGui())
	if autoStart {
		go StartServer()
	}
	return container.NewBorder(mainVerticalBox, nil, nil, nil)
}

func StartServer() {
	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	mux.HandleFunc("/slow", slowHandler)
	mux.HandleFunc("/query", queryParamsHandler)

	mux.HandleFunc("/pay", payHandler)
	mux.HandleFunc("/save", saveHandler)

	httpMtx.Lock()
	srv = &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}
	httpMtx.Unlock()

	fmt.Println("Starting server...")
	// ListenAndServe блокирующий
	// ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
	// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: http: Server closed» даже при нормальной остановке.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error:", err)
	}
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
