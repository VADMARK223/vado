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
	a := app.New()
	w := a.NewWindow("Vado")

	tabs := guiTabs.CreateAppTabs()

	// кнопка выхода
	exitBtn := gui.CreateBtn("Exit", theme.LogoutIcon(), func() { w.Close() })
	exitBtnWrapper := container.NewVBox(exitBtn)

	// верхняя панель = кнопка справа
	topBar := container.NewBorder(nil, nil, nil, exitBtnWrapper)

	// нижняя панель = версия справа
	bottomBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel(util.Tpl("Version %s", c.Version)),
	)

	// главный контейнер: сверху topBar, снизу bottomBar, центр = tabs
	root := container.NewBorder(topBar, bottomBar, nil, nil, tabs)
	w.SetContent(root)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.Close()
		}
	})

	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}
