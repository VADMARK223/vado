// Package gui содержит графические представления.
package gui

import (
	"fmt"
	"strings"
	c "vado/internal/constant"
	gui "vado/internal/gui/common"
	tabs "vado/internal/gui/tab"
	"vado/internal/gui/tab/settings"
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

func ShowMainApp(appCtx *util.AppContext) {
	logger.L().Info("Starting GUI mode.", zap.String("mode", util.GetModeValue()))
	a := app.NewWithID("io.vado")
	mainWindow := a.NewWindow("Vado")

	userInfo := widget.NewRichTextFromMarkdown(fmt.Sprintf("Пользователь: **%s**", "VADMARK"))
	modeTxt := widget.NewRichTextFromMarkdown(fmt.Sprintf("Режим: **%s**", strings.ToUpper(util.GetModeValue())))

	objs := []fyne.CanvasObject{userInfo, layout.NewSpacer()}
	objs = append(objs, createFastBlock()...)
	objs = append(objs,
		modeTxt,
		widget.NewRichTextFromMarkdown(fmt.Sprintf("Версия: **%s**", c.Version)),
		gui.NewBtn("", theme.LogoutIcon(), func() { mainWindow.Close() }),
	)

	bottomBar := container.NewHBox(objs...)
	root := container.NewBorder(nil, bottomBar, nil, nil, tabs.NewTabsView(appCtx, mainWindow))
	mainWindow.SetContent(root)

	mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			mainWindow.Close()
		}
	})

	mainWindow.Resize(fyne.NewSize(700, 400))
	mainWindow.ShowAndRun()
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
