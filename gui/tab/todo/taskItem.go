package todo

import (
	"fmt"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateTaskItem(task m.Task, deleteCallback func()) fyne.CanvasObject {
	taskLabel := widget.NewLabel(fmt.Sprintf("%d %s", task.Id, task.Name))
	hBox := container.NewHBox(taskLabel, layout.NewSpacer())
	taskDelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), deleteCallback)
	hBox.Add(taskDelBtn)
	return hBox
}
