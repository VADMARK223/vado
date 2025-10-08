package component

import (
	"strings"
	"time"
	t "vado/internal/domain/task"
	"vado/pkg/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ShowTaskDialog(parent fyne.Window, taskIn *t.Task, f func(t.Task)) {
	isEdit := taskIn != nil
	if taskIn == nil {
		taskIn = &t.Task{}
	}
	nameEntry := widget.NewEntry()
	descEntry := widget.NewMultiLineEntry()
	createAtEntry := widget.NewLabel("-")
	updateAtEntry := widget.NewLabel("-")
	check := widget.NewCheck("Выполнена", nil)
	var dlg dialog.Dialog
	saveBtn := widget.NewButton("", func() {
		taskId := func() int {
			if isEdit {
				return taskIn.ID
			}
			return 0
		}()
		outTask := t.Task{
			ID:          taskId,
			Name:        nameEntry.Text,
			Description: descEntry.Text,
			Completed:   check.Checked,
		}
		// Если режим редактирования, то прокидываем дату создания (потом переделать)
		if isEdit {
			outTask.CreatedAt = taskIn.CreatedAt
		}
		f(outTask)
		dlg.Hide()
	})
	saveBtn.Importance = widget.HighImportance
	saveBtn.Disable()

	var title string
	if isEdit {
		title = "Редактирование задачи"
		nameEntry.SetText(taskIn.Name)
		nameEntry.CursorColumn = len(taskIn.Name)
		descEntry.SetText(taskIn.Description)
		var temp string
		if taskIn.CreatedAt == nil {
			temp = "-"
		} else {
			temp = util.FormatTime(*taskIn.CreatedAt)
		}
		createAtEntry.SetText(temp)
		var updateAtEntryText string

		if taskIn.UpdatedAt == nil {
			updateAtEntryText = "-"
		} else {
			if *taskIn.CreatedAt == *taskIn.UpdatedAt {
				updateAtEntryText = "Не изменялась"
			} else {
				updateAtEntryText = util.FormatTime(*taskIn.UpdatedAt)
			}
		}

		updateAtEntry.SetText(updateAtEntryText)
		check.SetChecked(taskIn.Completed)
		saveBtn.SetText("Сохранить")
	} else {
		title = "Создание задачи"
		saveBtn.SetText("Создать")
	}

	cancelBtn := widget.NewButton("Отмена", func() {
		dlg.Hide()
	})

	nameEntry.OnChanged = func(text string) {
		updateSaveBtnEnable(saveBtn, text, taskIn.Name, descEntry.Text, taskIn.Description, check.Checked, taskIn.Completed)
	}

	descEntry.OnChanged = func(text string) {
		updateSaveBtnEnable(saveBtn, nameEntry.Text, taskIn.Name, text, taskIn.Description, check.Checked, taskIn.Completed)
	}

	check.OnChanged = func(checked bool) {
		updateSaveBtnEnable(saveBtn, nameEntry.Text, taskIn.Name, descEntry.Text, taskIn.Description, check.Checked, taskIn.Completed)
	}

	form := widget.NewForm(
		widget.NewFormItem("Название *", nameEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("", check),
		widget.NewFormItem("Создана", createAtEntry),
		widget.NewFormItem("Обновлена", updateAtEntry),
	)

	content := container.NewVBox(form, container.NewHBox(layout.NewSpacer(), cancelBtn, saveBtn))

	dlg = dialog.NewCustomWithoutButtons(title, content, parent)
	dlg.Resize(fyne.NewSize(400, 180))
	dlg.Show()

	// Через короткое время после показа диалога — установить фокус
	time.AfterFunc(100*time.Millisecond, func() {
		fyne.Do(func() {
			fyne.CurrentApp().Driver().CanvasForObject(nameEntry).Focus(nameEntry)
		})
	})
}

func updateSaveBtnEnable(btn *widget.Button, newName string, oldName string, newDesc string, oldDesc string, newCheck bool, oldCheck bool) {
	if getEnableSaveButton(newName, oldName, newDesc, oldDesc, newCheck, oldCheck) {
		btn.Enable()
	} else {
		btn.Disable()
	}
}

func getEnableSaveButton(newName string, oldName string, newDesc string, oldDesc string, newCheck bool, oldCheck bool) bool {
	if strings.TrimSpace(newName) == "" {
		return false
	}

	if newName == oldName && newDesc == oldDesc && newCheck == oldCheck {
		return false
	}

	return true
}
