package admin

import (
	"vado/internal/gui/common"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewUserControl() fyne.CanvasObject {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Введите пароль пользователя")

	createBtn := common.NewBtn("Создать", nil, func() {})

	grid := container.NewGridWithColumns(3,
		usernameEntry,
		passwordEntry,
		createBtn,
	)

	return grid
}
