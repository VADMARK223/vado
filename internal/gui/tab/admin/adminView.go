package admin

import (
	"vado/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewAdminView(appCtx *util.AppContext, win fyne.Window) fyne.CanvasObject {
	box := container.NewVBox(widget.NewLabel("Админка"), NewUserControl(appCtx, win))
	return box
}
