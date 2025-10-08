package admin

import (
	"context"
	"fmt"
	"image/color"
	"time"
	constant2 "vado/internal/constant"
	"vado/internal/gui/common"
	"vado/internal/gui/constant"
	"vado/internal/gui/tab/tasks/component"
	"vado/internal/pb/taskpb"
	appCtx "vado/internal/server/context"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewBoxGRPC(ctx *appCtx.AppContext, win fyne.Window) fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер GRPC:")

	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), func() {
		startServerGRPC(ctx)
	})
	startBtn.Disable()
	stopBtn := common.NewBtn("Стоп", theme.MediaStopIcon(), func() {
		stopServerGRPC(ctx)
	})
	statusIndicator := common.NewIndicator(color.RGBA{R: 255, G: 0, B: 0, A: 255}, fyne.NewSize(15, 15))

	go func() {
		for {
			running := ctx.GRPC.IsRunning()

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

	return container.NewHBox(lbl, startBtn, stopBtn, container.NewCenter(statusIndicator), NewClientBoxGRPC(win))
}

func startServerGRPC(ctx *appCtx.AppContext) {
	err := ctx.GRPC.Start()
	if err != nil {
		logger.L().Warn("Error start server grpc", zap.Error(err))
	}
}

func stopServerGRPC(ctx *appCtx.AppContext) {
	ctx.GRPC.Stop()
}

func NewClientBoxGRPC(win fyne.Window) fyne.CanvasObject {
	btn := common.NewBtn("Запросить кол-во заданий", nil, func() {
		conn, err := grpc.NewClient(fmt.Sprintf("localhost%s", constant2.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.L().Error("gRPC connection failed", zap.Error(err))
			dialog.ShowError(err, win)

		}
		defer func(conn *grpc.ClientConn) {
			if closeErr := conn.Close(); closeErr != nil {
				logger.L().Warn("Failed to close connection", zap.Error(closeErr))
			}
		}(conn)

		client := taskpb.NewTaskServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.GetAllTasks(ctx, nil)
		if err != nil {
			dialog.ShowError(err, win)
			logger.L().Error("gRPC request failed", zap.Error(err))
			return
		}

		dialog.ShowInformation("Ответ gRPC", fmt.Sprintf("Заданий: %d", len(resp.Tasks)), win)
	})

	return btn
}
