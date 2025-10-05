package component

import (
	"math/rand"
	m "vado/internal/model"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowEditTaskDialog(parent fyne.Window, task m.Task, f func(task m.Task)) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(task.Name)
	nameEntry.SetPlaceHolder("Редактирование задачи")
	nameEntry.Disable()

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(task.Description)
	descEntry.SetPlaceHolder("Описание задачи")
	descEntry.Disable()

	//var complete = task.Completed
	check := widget.NewCheck("Выполнено", func(checked bool) {
		logger.L().Debug("Change completed.")
		//complete = checked
	})
	check.SetChecked(task.Completed)
	check.Disable()

	formItems := []*widget.FormItem{
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("", check),
	}

	dlg := dialog.NewForm("Новая задача", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			logger.L().Debug("Save task.")
			updatedTask := m.Task{ID: rand.Intn(10000), Name: nameEntry.Text, Description: descEntry.Text, Completed: check.Checked}
			f(updatedTask)
		}
	}, parent)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()
}

/*func ShowEditTaskDialog(parent fyne.Window, task m.Task) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(task.Name)
	nameEntry.SetPlaceHolder("Редактирование задачи")
	nameEntry.Disable()

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetText(task.Description)
	descEntry.SetPlaceHolder("Описание задачи")
	descEntry.Disable()

	check := widget.NewCheck("Выполнено", func(checked bool) {
		logger.L().Debug("Change completed.")
	})
	check.Disable()

	formItems := []*widget.FormItem{
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("", check),
	}

	dlg := dialog.NewForm("Новая задача", "Сохранить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			logger.L().Debug("Save task.")
			//newTask := m.Task{ID: rand.Intn(10000), Name: nameEntry.Text, Description: descEntry.Text, Completed: complete}
			//f(newTask)
		}
	}, parent)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()
}*/
