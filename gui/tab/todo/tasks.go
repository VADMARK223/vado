package todo

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	c "vado/constant"
	"vado/gui/common"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateTODOTab(win fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Список заданий", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	data, err := os.ReadFile(c.TaskFileName)

	if err != nil {
		log.Fatal(err)
	}

	var list m.TaskList
	if err = json.Unmarshal(data, &list); err != nil {
		fmt.Printf("Error: '%s'\n", err)
	}
	vBox := container.NewVBox()

	refreshBtn := common.CreateBtn("Обновить", theme.ViewRefreshIcon(), func() {
		redrawList(list, vBox, win)
	})

	addBtn := common.CreateBtn("Добавить", theme.ContentAddIcon(), func() {
		showAddTaskDialog(win, &list, vBox)
	})
	controlBox := container.NewHBox(refreshBtn, addBtn)

	redrawList(list, vBox, win)

	scroll := container.NewVScroll(vBox)
	topBox := container.NewVBox(controlBox, title)
	content := container.NewBorder(topBox, nil, nil, nil, scroll)

	return content
}

// func redrawList(list m.TaskList, vBox *fyne.Container, win fyne.Window) {
func redrawList(list m.TaskList, vBox *fyne.Container, win fyne.Window) {
	vBox.RemoveAll()
	for i := range list.Tasks {
		task := list.Tasks[i]

		vBox.Add(CreateTaskItem(task, func() {
			//if common.ShowDeleteConfirm {
			//	dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить задание?", func(b bool) {
			//		if b {
			//			deleteTask(&list, task.Id)
			//			redrawList(list, vBox, win)
			//		}
			//	}, win)
			//} else {

			deleteTask(&list, task.Id)
			redrawList(list, vBox, win)

			//}
		}))
	}
}

func addTask(list *[]m.Task, task m.Task) {
	*list = append(*list, task)
}

func showAddTaskDialog(parent fyne.Window, list *m.TaskList, vBox *fyne.Container) {
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
			//desc := descEntry.Text

			newTask := m.Task{Id: rand.Intn(10000), Name: nameEntry.Text}
			addTask(&list.Tasks, newTask)
			saveJSON(list)

			var tempList []m.Task
			for _, t := range list.Tasks {
				tempList = append(tempList, t)
			}
			list.Tasks = tempList
			redrawList(*list, vBox, parent)
		}
	}, parent)

	dlg.Resize(fyne.NewSize(400, 300))
	dlg.Show()
}

func deleteTask(list *m.TaskList, taskId int) {
	var newTasks []m.Task
	for _, t := range list.Tasks {
		if t.Id != taskId {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks

	saveJSON(list)
}

func saveJSON(list *m.TaskList) {
	// Сохраняем в json
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		fmt.Println("Error marshal JSON:", err)
		return
	}
	/*
		0 → это префикс, который говорит Go, что число в восьмеричной системе.
		6 → владелец может читать и писать (4+2).
		4 → группа может только читать.
		4 → остальные пользователи могут только читать.
		os.ModePerm - это 0777, то есть маскимально открытые права.
	*/
	if err := os.WriteFile(c.TaskFileName, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
