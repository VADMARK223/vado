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

type TaskItem struct {
	widget.BaseWidget

	label    *widget.Label
	button   *widget.Button
	OnDelete func()
}

func NewTaskItem(text string, onDelete func()) *TaskItem {
	ti := &TaskItem{
		label:    widget.NewLabel(text),
		OnDelete: onDelete,
	}

	ti.button = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ti.OnDelete()
	})

	ti.ExtendBaseWidget(ti)
	return ti
}

func (ti *TaskItem) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewHBox(ti.label, layout.NewSpacer(), ti.button)
	return widget.NewSimpleRenderer(content)
}

// SetText обновляет текст задачи
func (ti *TaskItem) SetText(task m.Task) {
	text := util.Tpl("%d %s%s%s", task.Id, task.Name, func() string {
		if task.Description != "" {
			return util.Tpl(" (%s)", task.Description)
		}
		return ""
	}(), func() string {
		if task.Completed {
			return " Выполнено"
		}
		return " Не выполнено"
	}())

	ti.label.SetText(text)
}

// CreateTaskItem deprecated
func CreateTaskItem(task m.Task, deleteCallback func()) fyne.CanvasObject {
	text := util.Tpl("%d %s%s%s", task.Id, task.Name, func() string {
		if task.Description != "" {
			return util.Tpl(" (%s)", task.Description)
		}
		return ""
	}(), func() string {
		if task.Completed {
			return " Выполнено"
		}
		return " Не выполнено"
	}())
	taskLabel := widget.NewLabel(text)
	hBox := container.NewHBox(taskLabel, layout.NewSpacer())
	taskDelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), deleteCallback)
	hBox.Add(taskDelBtn)
	return hBox
}
