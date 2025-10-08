// Package gui содержит графические представления.
package gui

import (
	"fmt"
	"strings"
	c "vado/internal/constant"
	gui "vado/internal/gui/common"
	tabs "vado/internal/gui/tab"
	"vado/internal/gui/tab/settings"
	"vado/internal/server/context"
	"vado/internal/util"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

func ShowApp(ctx *context.AppContext) {
	logger.L().Info("Starting GUI mode.", zap.String("mode", util.GetModeValue()))
	a := app.NewWithID("io.vado")
	w := a.NewWindow("Vado")

	userInfo := widget.NewRichTextFromMarkdown(fmt.Sprintf("Пользователь: **%s**", "VADMARK"))
	modeTxt := widget.NewRichTextFromMarkdown(fmt.Sprintf("Режим: **%s**", strings.ToUpper(util.GetModeValue())))

	objs := []fyne.CanvasObject{userInfo, layout.NewSpacer()}
	objs = append(objs, createFastBlock()...)
	objs = append(objs,
		modeTxt,
		widget.NewRichTextFromMarkdown(fmt.Sprintf("Версия: **%s**", c.Version)),
		gui.NewBtn("", theme.LogoutIcon(), func() { w.Close() }),
	)

	bottomBar := container.NewHBox(objs...)
	root := container.NewBorder(nil, bottomBar, nil, nil, tabs.NewTabsView(ctx, w))
	w.SetContent(root)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			closeWindow(ctx, w)
		}
	})

	autoStartServers(ctx)

	w.SetCloseIntercept(func() {
		closeWindow(ctx, w)

	})
	w.Resize(fyne.NewSize(700, 400))
	w.ShowAndRun()
}

func closeWindow(ctx *context.AppContext, w fyne.Window) {
	ctx.HTTP.Stop()
	ctx.GRPC.Stop()
	w.Close()
}

func autoStartServers(ctx *context.AppContext) {
	if util.AutoStartServerHTTP() {
		err := ctx.HTTP.Start()
		if err != nil {
			logger.L().Error("HTTP server error", zap.Error(err))
			return
		}
	}

	if util.AutoStartServerGRPC() {
		if err := ctx.GRPC.Start(); err != nil {
			logger.L().Error("gRPC server error", zap.Error(err))
		}
	}

}

func createFastBlock() []fyne.CanvasObject {
	fastModeTxt := widget.NewRichTextFromMarkdown(getModeTxt(util.IsFastMode()))
	util.OnFastModeChange(func(newValue bool) {
		fastModeTxt.ParseMarkdown(getModeTxt(newValue))
	})

	if util.IsFastMode() {
		return []fyne.CanvasObject{
			fastModeTxt,
			settings.NewFastModeCheck(false),
		}
	}
	return nil
}

func getModeTxt(isMode bool) string {
	var mode string
	if isMode {
		mode = "ВКЛ"
	} else {
		mode = "ВЫКЛ"
	}

	return fmt.Sprintf("**%s**", mode)
}
