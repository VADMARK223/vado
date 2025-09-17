package http

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"sync"
	"time"
	"vado/constant"
	"vado/gui/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	httpMtx       sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func CreateHttpTab() fyne.CanvasObject {
	startBtn := common.CreateBtn("Start", theme.MediaPlayIcon(), StartServer)
	startBtn.Disable()

	stopBtn := common.CreateBtn("Stop", theme.MediaStopIcon(), stopServer)
	stopBtn.Disable()

	waitLbl := widget.NewLabel("Wait...")
	waitLbl.Hide()

	statusLbl := widget.NewLabel("Server status:")

	statusIndicator := canvas.NewCircle(color.White)
	statusIndicator.FillColor = constant.Red() // Red
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
					statusIndicator.FillColor = constant.Red()
				})
			} else {
				fyne.Do(func() {
					if inProcess {
						waitLbl.Show()
						startBtn.Disable()
						stopBtn.Disable()
						statusIndicator.FillColor = constant.Orange()
					} else {
						waitLbl.Hide()
						startBtn.Disable()
						stopBtn.Enable()
						statusIndicator.FillColor = constant.Green()
					}
				})
			}

			time.Sleep(time.Millisecond * constant.GuiUpdateMillisecond)
		}
	}()

	controlBox := container.NewHBox(startBtn, stopBtn, waitLbl)
	mainVerticalBox := container.NewVBox(container.NewHBox(statusLbl, statusIndicatorLayout), controlBox)
	mainVerticalBox.Add(widget.NewSeparator())
	mainVerticalBox.Add(createMoneyGui())
	if constant.AutoStart {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*constant.ShutdownSecond)
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
