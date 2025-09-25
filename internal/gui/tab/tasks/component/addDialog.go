package component

import (
	"math/rand"
	m "vado/internal/model"
	"vado/pkg/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowAddTaskDialog(parent fyne.Window, f func(task m.Task)) {
	nameEntry := widget.NewEntry()
	initName := util.Tpl("Задача %d", rand.Intn(10))
	nameEntry.SetText(initName)
	nameEntry.SetPlaceHolder("Название задачи")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Описание задачи")

	var complete = false
	check := widget.NewCheck("Выполнено", func(checked bool) {
		complete = checked
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("", check),
	}
	dlg := dialog.NewForm("Новая задача", "Добавить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			newTask := m.Task{ID: rand.Intn(10000), Name: nameEntry.Text, Description: descEntry.Text, Completed: complete}
			f(newTask)
		}
	}, parent)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()
}
