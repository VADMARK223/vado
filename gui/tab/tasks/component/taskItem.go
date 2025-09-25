package component

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
	widget.BaseWidget // Встраивание

	check    *widget.Check
	label    *widget.Label
	button   *widget.Button
	OnDelete func()
}

func NewTaskItem(text string, onDelete func()) *TaskItem {
	check := widget.NewCheck("", func(b bool) {})
	check.Disable()

	ti := &TaskItem{
		check:    check,
		label:    widget.NewLabel(text),
		OnDelete: onDelete,
	}

	ti.button = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ti.OnDelete()
	})

	ti.ExtendBaseWidget(ti) // Сообщаем движку, что не простая структура
	return ti
}

func (ti *TaskItem) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewHBox(ti.check, ti.label, layout.NewSpacer(), ti.button)
	return widget.NewSimpleRenderer(content)
}

// SetTask прокидывает данные из представления
func (ti *TaskItem) SetTask(task m.Task) {
	ti.check.SetChecked(task.Completed)

	text := util.Tpl("%d %s%s", task.ID, task.Name, GetDescText(task.Description))

	ti.label.SetText(text)
}

func GetDescText(taskDesc string) string {
	if taskDesc != "" {
		return util.Tpl(" (%s)", taskDesc)
	}
	return ""
}
