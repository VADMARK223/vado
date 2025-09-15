package tabs

import (
	"context"
	"fmt"
	"image/color"
	"net/http"
	"sync"
	"time"
	"vado/gui/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	srv           *http.Server // Глобальный сервер, один на все вызовы
	mtx           sync.Mutex
	stopInProcess = false // Сделать потокобезопастным
)

func CreateHttpTab() fyne.CanvasObject {
	startBtn := common.CreateBtn("Start", nil)
	startBtn.Disable()

	stopBtn := common.CreateBtn("Stop", nil)
	stopBtn.Disable()

	waitLbl := widget.NewLabel("Wait...")
	waitLbl.Hide()

	statusLbl := widget.NewLabel("Server status:")

	statusIndicator := canvas.NewCircle(color.White)
	statusIndicator.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red
	statusIndicator.StrokeColor = color.Gray{Y: 0x99}
	statusIndicator.StrokeWidth = 1
	statusIndicatorLayout := container.NewWithoutLayout(statusIndicator)
	statusIndicator.Resize(fyne.NewSize(30, 30))

	startBtn.OnTapped = func() {
		go startServer()
	}

	stopBtn.OnTapped = func() {
		stopInProcess = true
		go stopServer()
	}

	go func() {
		for {
			mtx.Lock()
			running := srv == nil
			mtx.Unlock()
			if running {
				fyne.Do(func() {
					waitLbl.Hide()
					startBtn.Enable()
					stopBtn.Disable()
					statusIndicator.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				})
			} else {
				fyne.Do(func() {
					if stopInProcess {
						waitLbl.Show()
						startBtn.Disable()
						stopBtn.Disable()
						statusIndicator.FillColor = color.RGBA{R: 255, G: 165, B: 0, A: 255}
					} else {
						waitLbl.Hide()
						startBtn.Disable()
						stopBtn.Enable()
						statusIndicator.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
					}
				})
			}

			time.Sleep(time.Millisecond * 500)
		}
	}()

	hBox := container.NewHBox(startBtn, stopBtn, waitLbl)
	vBox := container.NewVBox(container.NewHBox(statusLbl, statusIndicatorLayout), hBox)
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
	// ListenAndServe блокирующий → запускаем в горутине
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error:", err)
	}
}

func slowHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Started slow request...")
	time.Sleep(time.Second * 10)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
