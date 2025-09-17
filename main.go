package main

import (
	c "vado/constant"
	gui "vado/gui/common"
	guiTabs "vado/gui/tabs"
	"vado/gui/tabs/http"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	if c.ShowGui {
		showGui()
	} else {
		http.StartServer()
	}
}

func showGui() {
	a := app.NewWithID("io.vado")
	mainWindow := a.NewWindow("Vado")

	tabs := guiTabs.CreateAppTabs(mainWindow)
	exitBtn := gui.CreateBtn("", theme.LogoutIcon(), func() { mainWindow.Close() })
	topBar := container.NewBorder(nil, nil, nil, exitBtn)

	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(util.Tpl("Version %s", c.Version)),
	)

	root := container.NewBorder(topBar, bottomBar, nil, nil, tabs)
	mainWindow.SetContent(root)

	mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			mainWindow.Close()
		}
	})

	mainWindow.Resize(fyne.NewSize(350, 400))
	mainWindow.ShowAndRun()
}
