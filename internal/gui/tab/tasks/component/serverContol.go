package component

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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	guiUpdateMillisecond   = 500 // Частота обновления GUI
	shutdownSecond         = 5
	slowRequestDelaySecond = 10 // Длительность выполнения медленного запроса// Время "мягкой" остановки сервер
)

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	httpMtx       sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func NewServerControl() fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер:")
	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), nil)
	startBtn.OnTapped = func() {
		go func() {
			if err := startServer(); err != nil {
				fmt.Println("Start server error:", err)
			}
		}()
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
	//statusIndicator := canvas.NewCircle(color.White)
	//statusIndicator.FillColor = constant2.Red()
	//statusIndicator.StrokeColor = color.Gray{Y: 0x99}
	//statusIndicator.StrokeWidth = 1
	//statusIndicator.Resize(fyne.NewSize(15, 15))
	//statusIndicatorLayout := container.NewWithoutLayout(statusIndicator)

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

	return container.NewHBox(lbl, startBtn, stopBtn, container.NewCenter(statusIndicator), waitLbl)
}

func startServer() error {
	httpMtx.Lock()
	if srv != nil {
		return errors.New("server already running")
	}

	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	mux.HandleFunc("/slow", slowHandler)
	mux.HandleFunc("/tasks", taskHandler)

	srv = &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}
	httpMtx.Unlock()

	go func() {
		fmt.Println("Tasks server started on :9091...")

		// ListenAndServe блокирующий
		// ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
		// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: http: Server closed» даже при нормальной остановке.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Error:", err)
		}
	}()

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
