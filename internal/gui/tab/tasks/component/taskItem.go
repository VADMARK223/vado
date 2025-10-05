package component

import (
	"time"
	m "vado/internal/model"
	"vado/pkg/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TaskItem struct {
	widget.BaseWidget

	check         *widget.Check
	label         *widget.Label
	createAtLabel *widget.Label
	button        *widget.Button
	OnDelete      func()
	OnDoubleTap   func()

	lastTap time.Time
}

func NewTaskItem(text string, onDelete func()) *TaskItem {
	check := widget.NewCheck("", func(b bool) {})
	check.Disable()

	ti := &TaskItem{
		check:         check,
		label:         widget.NewLabel(text),
		createAtLabel: widget.NewLabel(""),
		OnDelete:      onDelete,
	}

	ti.button = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		ti.OnDelete()
	})

	ti.ExtendBaseWidget(ti) // Сообщаем движку, что не простая структура
	return ti
}

func (ti *TaskItem) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewHBox(ti.check, ti.label, layout.NewSpacer(), ti.createAtLabel, ti.button)
	return widget.NewSimpleRenderer(content)
}

func (ti *TaskItem) Tapped(ev *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(ti.lastTap) < 300*time.Millisecond {
		ti.lastTap = time.Time{}
		if ti.OnDoubleTap != nil {
			ti.OnDoubleTap()
		}
	} else {
		ti.lastTap = now
	}
}

// SetTask прокидывает данные из представления
func (ti *TaskItem) SetTask(task m.Task) {
	ti.check.SetChecked(task.Completed)

	text := util.Tpl("%d %s%s", task.ID, task.Name, GetDescText(task.Description))
	ti.label.SetText(text)

	var value string
	if task.CreatedAt != nil {
		value = util.FormatTime(*task.CreatedAt)
	} else {
		value = "-"
	}
	ti.createAtLabel.SetText(value)
	createdAtText := util.Tpl("Создана: %s", value)
	ti.createAtLabel.SetText(createdAtText)
}

func GetDescText(taskDesc string) string {
	if taskDesc != "" {
		return util.Tpl(" (%s)", taskDesc)
	}
	return ""
}
