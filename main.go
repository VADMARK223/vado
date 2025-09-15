package main

import (
	c "vado/constant"
	gui "vado/gui/common"
	guiTabs "vado/gui/tabs"
	"vado/module/http"
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
	a := app.New()
	w := a.NewWindow("Vado")

	tabs := guiTabs.CreateAppTabs()

	exitBtn := gui.CreateBtn("Exit", theme.LogoutIcon(), func() { w.Close() })
	exitBtnWrapper := container.NewVBox(exitBtn)
	topBar := container.NewBorder(nil, nil, nil, exitBtnWrapper, tabs)
	bottomBar := container.NewHBox(layout.NewSpacer(), widget.NewLabel(util.Tpl("Version %s", c.Version)))
	header := container.NewBorder(topBar, bottomBar, nil, nil)
	w.SetContent(header)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.Close()
		}
	})
	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
