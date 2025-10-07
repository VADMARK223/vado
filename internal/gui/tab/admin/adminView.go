package admin

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewAdminView() fyne.CanvasObject {
	box := container.NewVBox(widget.NewLabel("Админка"), NewUserControl())
	return box
}
