package conponent

import (
	"fmt"
	"math/rand"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowAddTaskDialog(parent fyne.Window, f func(task m.Task)) {
	nameEntry := widget.NewEntry()
	initName := fmt.Sprintf("Задача %d", rand.Intn(10))
	nameEntry.SetText(initName)
	nameEntry.SetPlaceHolder("Название задачи")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Описание задачи")

	formItems := []*widget.FormItem{
		widget.NewFormItem("Название", nameEntry),
		widget.NewFormItem("Описание", descEntry),
	}
	dlg := dialog.NewForm("Новая задача", "Добавить", "Отмена", formItems, func(confirm bool) {
		if confirm {
			newTask := m.Task{Id: rand.Intn(10000), Name: nameEntry.Text, Desc: descEntry.Text}
			f(newTask)
		}
	}, parent)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()
}
