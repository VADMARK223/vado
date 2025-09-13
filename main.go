package main

import (
	c "vado/constant"
	"vado/gui/modules"
	"vado/gui/tasks"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Vado")

	tabs := container.NewAppTabs(
		container.NewTabItem("Modules", modules.CreateModulesGui()),
		container.NewTabItem("Tasks", tasks.CreateTasksGui()))
	tabs.SetTabLocation(container.TabLocationTop)

	exitBtn := widget.NewButton("Exit", func() { w.Close() })
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
