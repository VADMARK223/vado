package tabs

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateSettingsTab() fyne.CanvasObject {
	label := widget.NewLabel("In development")
	time.Sleep(2 * time.Second)
	return label
}
