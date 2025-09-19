package conponent

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func CreateTasksTitle() *canvas.Text {
	title := canvas.NewText("Список заданий", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	return title
}
