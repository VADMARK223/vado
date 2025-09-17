package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "fyne.io/fyne/v2/widget"
)

func CreateLazyTabItem(text string, factory func() fyne.CanvasObject, factories map[*container.TabItem]func() fyne.CanvasObject) *container.TabItem {
	tab := container.NewTabItem(text, widget.NewLabel("..."))
	factories[tab] = factory
	return tab
}
