package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/widget"
)

func CreateTabItem(text string, content fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(text, content)
}
