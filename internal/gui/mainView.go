// Package gui содержит графические представления.
package gui

import (
	"fmt"
	c "vado/internal/constant"
	gui "vado/internal/gui/common"
	tabs "vado/internal/gui/tab"
	"vado/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ShowMainApp() {
	a := app.NewWithID("io.vado")
	mainWindow := a.NewWindow("Vado")

	modeTxt := widget.NewRichTextFromMarkdown(getModeTxt(util.IsDevMode()))
	util.OnDevModeChange(func(newValue bool) {
		modeTxt.ParseMarkdown(getModeTxt(newValue))
	})

	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		modeTxt,
		widget.NewRichTextFromMarkdown(fmt.Sprintf("Версия: **%s**", c.Version)),
		gui.NewBtn("", theme.LogoutIcon(), func() { mainWindow.Close() }),
	)

	root := container.NewBorder(nil /*topBar*/, bottomBar, nil, nil, tabs.NewTabsView(mainWindow))
	mainWindow.SetContent(root)

	mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			mainWindow.Close()
		}
	})

	mainWindow.Resize(fyne.NewSize(700, 400))
	mainWindow.ShowAndRun()
}

func getModeTxt(isMode bool) string {
	var mode string
	if isMode {
		mode = "DEV"
	} else {
		mode = "PROD"
	}

	return fmt.Sprintf("Режим: **%s**", mode)
}
