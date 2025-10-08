package admin

import (
	"image/color"
	"time"
	"vado/internal/gui/common"
	"vado/internal/gui/constant"
	"vado/internal/gui/tab/tasks/component"
	appCtx "vado/internal/server/context"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

func NewBoxHTTP(ctx *appCtx.AppContext, win fyne.Window) fyne.CanvasObject {
	lbl := widget.NewLabel("Сервер HTTP:")

	startBtn := common.NewBtn("Старт", theme.MediaPlayIcon(), func() {
		startServerHTTP(ctx)
	})
	startBtn.Disable()

	stopBtn := common.NewBtn("Стоп", theme.MediaStopIcon(), func() {
		stopServerHTTP(ctx)
	})

	statusIndicator := common.NewIndicator(color.RGBA{R: 255, G: 0, B: 0, A: 255}, fyne.NewSize(15, 15))

	// Фоновое обновление статуса
	go func() {
		for {
			running := ctx.HTTP.IsRunning()

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

	return container.NewHBox(
		lbl,
		startBtn,
		stopBtn,
		container.NewCenter(statusIndicator),
	)
}

func startServerHTTP(ctx *appCtx.AppContext) {
	err := ctx.HTTP.Start()
	if err != nil {
		logger.L().Warn("Error start server http", zap.Error(err))
	}
}

func stopServerHTTP(ctx *appCtx.AppContext) {
	ctx.HTTP.Stop()
}
