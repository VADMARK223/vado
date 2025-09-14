package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateSettingsGui() fyne.CanvasObject {
	label := widget.NewLabel("In development")
	return label
}
