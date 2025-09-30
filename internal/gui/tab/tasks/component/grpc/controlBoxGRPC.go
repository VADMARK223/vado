package grpc

import (
	"fmt"
	"image/color"
	"log"
	"net"
	"sync"
	"time"
	"vado/internal/gui/common"
	"vado/internal/gui/constant"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/pb/taskpb"
	"vado/internal/service"
	"vado/internal/util"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServer *grpc.Server
	mu         sync.Mutex
	listener   net.Listener
)

func NewControlBoxGRPC(s service.ITaskService, win fyne.Window) fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер GRPC:")

	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), func() {
		startServerGRPC(s)
	})
	startBtn.Disable()
	stopBtn := common.NewBtn("Стоп", theme.MediaStopIcon(), func() {
		stopServerGRPC()
	})
	statusIndicator := common.NewIndicator(color.RGBA{R: 255, G: 0, B: 0, A: 255}, fyne.NewSize(15, 15))

	go func() {
		for {
			mu.Lock()
			running := grpcServer != nil
			mu.Unlock()

			fyne.Do(func() {
				if running {
					startBtn.Disable()
					stopBtn.Enable()
					statusIndicator.SetFillColor(constant.Green())
				} else {
					startBtn.Enable()
					stopBtn.Disable()
					statusIndicator.SetFillColor(constant.Red())
				}
			})
			time.Sleep(time.Millisecond * component.GuiUpdateMillisecond)
		}
	}()

	if util.AutoStartServerGRPC() {
		startServerGRPC(s)
	}

	return container.NewHBox(lbl, startBtn, stopBtn, container.NewCenter(statusIndicator), NewClientBoxGRPC(win))
}

func startServerGRPC(s service.ITaskService) {
	mu.Lock()
	defer mu.Unlock()

	if grpcServer != nil {
		fmt.Println("Server not running")
		return
	}

	var err error
	listener, err = net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	grpcServer = grpc.NewServer()
	serviceGRPC := service.NewTaskServiceGRPC(s)
	taskpb.RegisterTaskServiceServer(grpcServer, serviceGRPC)
	reflection.Register(grpcServer)

	go func() {
		logger.L().Info("gRPC-server started on :50051")

		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()
}

func stopServerGRPC() {
	mu.Lock()
	defer mu.Unlock()

	if grpcServer == nil {
		fmt.Println("Server not running")
		return
	}

	grpcServer.GracefulStop()
	grpcServer = nil
	listener = nil

	fmt.Println("GRPC server stopped.")
}
