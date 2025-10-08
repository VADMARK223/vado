package admin

import (
	"vado/internal/server/context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewAdminView(appCtx *context.AppContext, win fyne.Window) fyne.CanvasObject {
	grpc := NewBoxGRPC(appCtx, win)
	http := NewBoxHTTP(appCtx, win)
	box := container.NewVBox(http, grpc, NewUserControl(appCtx, win))
	return box
}
