package heavy

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewHeavyView() fyne.CanvasObject {
	label := widget.NewLabel("В разработке...")
	time.Sleep(2 * time.Second)
	return label
}
