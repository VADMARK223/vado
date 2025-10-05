package component

import (
	"strings"
	m "vado/internal/model"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowTaskDialog(parent fyne.Window, task *m.Task, f func(task m.Task)) {
	isEdit := task != nil
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	check := widget.NewCheck("Выполнена", nil)
	var dlg dialog.Dialog
	saveBtn := widget.NewButton("", func() {
		logger.L().Debug("Save task.")
		taskId := func() int {
			if isEdit {
				return task.ID
			}
			return -1
		}()
		updatedTask := m.Task{
			ID:          taskId,
			Name:        nameEntry.Text,
			Description: descEntry.Text,
			Completed:   check.Checked,
		}
		f(updatedTask)
		dlg.Hide()
	})
	saveBtn.Importance = widget.HighImportance
	saveBtn.Disable()

	var title string
	if isEdit {
		title = "Редактирование задачи"
		nameEntry.SetText(task.Name)
		descEntry.SetText(task.Description)
		check.SetChecked(task.Completed)
		saveBtn.SetText("Сохранить")
		saveBtn.Enable()
	} else {
		title = "Создание задачи"
		saveBtn.SetText("Создать")
	}

	cancelBtn := widget.NewButton("Отмена", func() {
		dlg.Hide()
	})

	nameEntry.OnChanged = func(text string) {
		saveBtn.Enable()
		if strings.TrimSpace(text) == "" {
			saveBtn.Disable()
		}
	}

	form := widget.NewForm(
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("", check),
	)

	content := container.NewVBox(form, container.NewHBox(layout.NewSpacer(), cancelBtn, saveBtn))

	dlg = dialog.NewCustomWithoutButtons(title, content, parent)
	dlg.Resize(fyne.NewSize(400, 180))
	dlg.Show()
}
