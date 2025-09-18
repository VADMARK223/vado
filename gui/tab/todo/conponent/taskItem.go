package conponent

import (
	m "vado/model"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateTaskItem(task m.Task, deleteCallback func()) fyne.CanvasObject {
	text := util.Tpl("%d %s%s", task.Id, task.Name, func() string {
		if task.Desc != "" {
			return util.Tpl(" (%s)", task.Desc)
		}
		return ""
	}())
	taskLabel := widget.NewLabel(text)
	hBox := container.NewHBox(taskLabel, layout.NewSpacer())
	taskDelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), deleteCallback)
	hBox.Add(taskDelBtn)
	return hBox
}
