package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateSettingsTab() fyne.CanvasObject {
	label := widget.NewLabel("In development")
	return label
}
