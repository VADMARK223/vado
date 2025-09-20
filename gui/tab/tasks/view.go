package tasks

import (
	"image/color"
	"vado/gui/common"
	"vado/gui/tab/tasks/conponent"
	"vado/gui/tab/tasks/constant"
	m "vado/model"
	"vado/service"
	"vado/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ViewTasks struct {
	service     *service.TaskService
	untypedList binding.UntypedList
}

func NewTasksView(win fyne.Window, s *service.TaskService) fyne.CanvasObject {
	vt := &ViewTasks{service: s, untypedList: binding.NewUntypedList()}
	err := vt.reloadTasks()
	if err != nil {
		return nil
	}

	refreshBtn := common.CreateBtn("Обновить", theme.ViewRefreshIcon(), nil)
	refreshBtn.Disable()
	addBtn := common.CreateBtn("Добавить", theme.ContentAddIcon(), func() {
		conponent.ShowAddTaskDialog(win, func(newTask m.Task) {
			vt.AddTask(newTask)
		})
	})
	quickAddBtn := common.CreateBtn("Добавить (быстро)", theme.ContentAddIcon(), func() {
		vt.AddTaskQuick()
	})

	deleteAllBtn := common.CreateBtn("Удалить все", theme.DeleteIcon(), func() {
		if constant.ShowTaskDeleteALLConfirm {
			dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить все задания?", func(b bool) {
				if b {
					vt.DeleteAllTasks()
				}
			}, win)
		} else {
			vt.DeleteAllTasks()
		}
	})

	list := widget.NewListWithData(
		vt.untypedList,
		func() fyne.CanvasObject {
			return conponent.NewTaskItem("", nil)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			task, _ := item.(binding.Untyped).Get()
			t := task.(m.Task)

			taskItem := obj.(*conponent.TaskItem)
			taskItem.SetTask(t)

			doDelete := func() {
				delErr := vt.service.Delete(t.Id)
				if delErr != nil {
					panic(delErr)
					return
				}
				_ = vt.reloadTasks()
				if err != nil {
					panic(err)
					return
				}
			}

			taskItem.OnDelete = func() {
				if constant.ShowTaskDeleteConfirm {
					dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить задание?", func(b bool) {
						if b {
							doDelete()
						}
					}, win)
				} else {
					doDelete()
				}
			}
		})
	scroll := container.NewVScroll(list)
	controlBox := container.NewHBox(refreshBtn, addBtn, quickAddBtn, layout.NewSpacer(), deleteAllBtn)
	title := canvas.NewText("Список заданий", color.White)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	vt.untypedList.AddListener(binding.NewDataListener(func() {
		tasksListLen := vt.untypedList.Length()
		if tasksListLen == 0 {
			deleteAllBtn.Disable()
			title.Text = "Список заданий пуст."
		} else {
			deleteAllBtn.Enable()
			title.Text = util.Tpl("Список заданий (Всего: %d)", tasksListLen)
		}
	}))
	topBox := container.NewVBox(controlBox, title)
	content := container.NewBorder(topBox, nil, nil, nil, scroll)
	return content

}

func (vt *ViewTasks) reloadTasks() error {
	tasksList, err := vt.service.GetAllTasks()
	if err != nil {
		return err
	}

	// преобразуем []Task → []any, потому что UntypedList принимает any
	items := make([]any, len(tasksList.Tasks))
	for i, t := range tasksList.Tasks {
		items[i] = t
	}

	return vt.untypedList.Set(items)
}
