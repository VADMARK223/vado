package tasks

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"vado/gui/common"
	"vado/gui/tab/tasks/conponent"
	"vado/gui/tab/tasks/constant"
	m "vado/model"
	"vado/service"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ViewTasks struct {
	service *service.TaskService
	list    *widget.List
	tasks   []m.Task
}

const newVersion = true

func CreateView(win fyne.Window, s *service.TaskService) fyne.CanvasObject {
	vt := &ViewTasks{service: s}
	vt.reloadTasks()
	tasksList, _ := s.GetAllTasks()

	vBox := container.NewVBox()

	refreshBtn := common.CreateBtn("Обновить", theme.ViewRefreshIcon(), func() {
		redraw(&tasksList, vBox, win)
	})

	addBtn := common.CreateBtn("Добавить", theme.ContentAddIcon(), func() {
		conponent.ShowAddTaskDialog(win, func(task m.Task) {
			addSaveRedraw(&tasksList, vBox, win, task)
		})
	})

	quickAddBtn := common.CreateBtn("Добавить (быстро)", theme.ContentAddIcon(), func() {
		if newVersion {
			vt.AddTaskFast("Fast task")
		} else {
			id := rand.Intn(10000)
			addSaveRedraw(&tasksList, vBox, win, m.Task{Id: id, Name: util.Tpl("Задание %d", id)})
		}
	})

	deleteAllBtn := common.CreateBtn("Удалить все", theme.DeleteIcon(), func() {
		tasksList.Tasks = []m.Task{}
		saveJSON(&tasksList)
		redraw(&tasksList, vBox, win)
	})

	vt.list = widget.NewList(func() int {
		return len(tasksList.Tasks)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		item := object.(*widget.Label)
		item.SetText(tasksList.Tasks[id].Name)
	})
	scroll := container.NewVScroll(vt.list)
	controlBox := container.NewHBox(refreshBtn, addBtn, quickAddBtn, layout.NewSpacer(), deleteAllBtn)
	topBox := container.NewVBox(controlBox, conponent.CreateTasksTitle())
	content := container.NewBorder(topBox, nil, nil, nil, scroll)
	return content

}

func (vt *ViewTasks) reloadTasks() {
	tasksList, _ := vt.service.GetAllTasks()
	vt.tasks = tasksList.Tasks
}

func addSaveRedraw(list *m.TaskList, listBox *fyne.Container, win fyne.Window, task m.Task) {
	if err := addTask(&list.Tasks, task, win); err != nil {
		fmt.Printf("Error: %s id = %d\n", err, task.Id)
		return
	}
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

func addTask(tasks *[]m.Task, newTask m.Task, win fyne.Window) error {
	for _, t := range *tasks {
		if t.Id == newTask.Id {
			dialog.NewCustom("Ошибка", "OK", widget.NewLabel("Задача с таким ID уже существует"), win).Show()
			return constant.ErrTaskAlreadyExists
		}
	}

	*tasks = append(*tasks, newTask)

	return nil
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
		Флаг os. ModePerm - это 0777, то есть максимально открытые права.
	*/
	if err := os.WriteFile(constant.TaskFilePath, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
