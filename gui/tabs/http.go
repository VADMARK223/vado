package tabs

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
	mtx           sync.Mutex
	stopInProcess = false // Сервер в процессе остановки
)

func CreateHttpTab() fyne.CanvasObject {
	startBtn := common.CreateBtn("Start", theme.MediaPlayIcon(), startServer)
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
		go startServer()
	}

	stopBtn.OnTapped = func() {
		mtx.Lock()
		if srv == nil || stopInProcess {
			mtx.Unlock()
			return
		}
		stopInProcess = true
		mtx.Unlock()
		go stopServer()
	}

	go func() {
		for {
			mtx.Lock()
			running := srv == nil
			inProcess := stopInProcess
			mtx.Unlock()
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

	hBox := container.NewHBox(startBtn, stopBtn, waitLbl)
	vBox := container.NewVBox(container.NewHBox(statusLbl, statusIndicatorLayout), hBox)
	if constant.AutoStart {
		go startServer()
	}
	return container.NewBorder(vBox, nil, nil, nil)
}

func startServer() {
	mux := http.NewServeMux() // multiplexer = «распределитель запросов»
	mux.HandleFunc("/slow", slowHandler)

	mtx.Lock()
	srv = &http.Server{
		Addr:    ":9091",
		Handler: mux,
	}
	mtx.Unlock()

	fmt.Println("Starting server...")
	// ListenAndServe блокирующий
	// http.ErrServerClosed это не ошибка, а сигнал: «Сервер завершён штатно».
	// Поэтому её нужно отфильтровать, иначе в логах всегда будет «Error: http: Server closed» даже при нормальном стопе.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Error:", err)
	}
}

func slowHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Started slow request...")
	time.Sleep(time.Second * constant.SlowRequestDelaySecond)
	str := "Hello from slow handler!"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Finished slow request")
	}
}

func stopServer() {
	mtx.Lock()
	s := srv
	if s == nil {
		panic("Server is nil!")
	}
	mtx.Unlock()

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
	mtx.Lock()
	stopInProcess = false
	srv = nil
	mtx.Unlock()
}
