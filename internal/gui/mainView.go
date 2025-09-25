// Package gui содержит графические представления.
package gui

import (
	c "vado/internal/constant"
	gui "vado/internal/gui/common"
	tabs "vado/internal/gui/tab"
	"vado/pkg/util"

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

	exitBtn := gui.CreateBtn("", theme.LogoutIcon(), func() { mainWindow.Close() })
	topBar := container.NewBorder(nil, nil, nil, exitBtn)

	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(util.Tpl("Version %s", c.Version)),
	)

	root := container.NewBorder(topBar, bottomBar, nil, nil, tabs.NewTabsView(mainWindow))
	mainWindow.SetContent(root)

	mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			mainWindow.Close()
		}
	})

	mainWindow.Resize(fyne.NewSize(700, 400))
	mainWindow.ShowAndRun()
}
