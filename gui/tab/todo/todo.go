package todo

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"vado/gui/common"
	"vado/gui/tab/todo/conponent"
	"vado/gui/tab/todo/constant"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateTODOTab(win fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Список заданий", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	data, err := os.ReadFile(constant.TaskFileName)

	if err != nil {
		log.Fatal(err)
	}

	var list m.TaskList
	if err = json.Unmarshal(data, &list); err != nil {
		fmt.Printf("Error: '%s'\n", err)
	}
	vBox := container.NewVBox()

	refreshBtn := common.CreateBtn("Обновить", theme.ViewRefreshIcon(), func() {
		redraw(&list, vBox, win)
	})

	addBtn := common.CreateBtn("Добавить", theme.ContentAddIcon(), func() {
		conponent.ShowAddTaskDialog(win, func(task m.Task) {
			addSaveRedraw(&list, vBox, win, task)
		})
	})

	quickAddBtn := common.CreateBtn("Добавить (быстро)", theme.ContentAddIcon(), func() {
		id := rand.Intn(1000)
		addSaveRedraw(&list, vBox, win, m.Task{Id: id, Name: fmt.Sprintf("Задаиние %d", id)})
	})

	deleteAllBtn := common.CreateBtn("Удалить все", theme.DeleteIcon(), func() {
		list.Tasks = nil
		saveJSON(&list)
		redraw(&list, vBox, win)
	})
	findInt := widget.NewEntry()
	findInt.SetPlaceHolder("Поиск")
	findInt.Hide()

	controlBox := container.NewHBox(refreshBtn, addBtn, quickAddBtn, layout.NewSpacer(), findInt, deleteAllBtn)

	redraw(&list, vBox, win)

	scroll := container.NewVScroll(vBox)
	topBox := container.NewVBox(controlBox, title)
	content := container.NewBorder(topBox, nil, nil, nil, scroll)

	return content
}

func addSaveRedraw(list *m.TaskList, listBox *fyne.Container, win fyne.Window, task m.Task) {
	addTask(&list.Tasks, task)
	saveJSON(list)
	redraw(list, listBox, win)
}

func redraw(list *m.TaskList, listBox *fyne.Container, win fyne.Window) {
	listBox.RemoveAll()
	for i := range list.Tasks {
		task := list.Tasks[i]

		listBox.Add(conponent.CreateTaskItem(task, func() {
			doDelete := func() {
				deleteTaskAndRedraw(list, task.Id, listBox, win)
			}
			if constant.ShowDeleteConfirm {
				dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить задание?", func(b bool) {
					if b {
						doDelete()
					}
				}, win)
			} else {
				doDelete()
			}
		}))
	}
}

func deleteTaskAndRedraw(list *m.TaskList, taskId int, listBox *fyne.Container, win fyne.Window) {
	deleteTask(list, taskId)
	redraw(list, listBox, win)
}

func addTask(list *[]m.Task, task m.Task) {
	*list = append(*list, task)
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
		os.ModePerm - это 0777, то есть максимально открытые права.
	*/
	if err := os.WriteFile(constant.TaskFileName, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
