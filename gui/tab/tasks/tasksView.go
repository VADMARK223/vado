package tasks

import (
	"fmt"
	"image/color"
	constant2 "vado/constant"
	"vado/gui/common"
	"vado/gui/tab/tasks/component"
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
	service     service.ITaskService
	untypedList binding.UntypedList
}

func NewTasksView(win fyne.Window, s service.ITaskService, isJSON bool) fyne.CanvasObject {
	vt := &ViewTasks{service: s, untypedList: binding.NewUntypedList()}
	err := vt.reloadTasks()
	if err != nil {
		return nil
	}

	modeLbl := widget.NewRichTextFromMarkdown(func() string {
		var mode string
		if isJSON {
			mode = "JSON"
		} else {
			mode = "DB"
		}

		return fmt.Sprintf("Режим: **%s**", mode)
	}())

	refreshBtn := common.CreateBtn("", theme.ViewRefreshIcon(), nil)
	refreshBtn.Disable()
	addBtn := common.CreateBtn("Добавить", theme.ContentAddIcon(), func() {
		component.ShowAddTaskDialog(win, func(newTask m.Task) {
			vt.AddTask(newTask)
		})
	})
	quickAddBtn := common.CreateBtn("Добавить (быстро)", theme.ContentAddIcon(), func() {
		vt.AddTaskQuick()
	})

	deleteAllBtn := common.CreateBtn("Удалить все", theme.DeleteIcon(), func() {
		if util.GetBoolPrefByKey(constant2.DevModePref) {
			vt.DeleteAllTasks()
		} else {
			dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить все задания?", func(b bool) {
				if b {
					vt.DeleteAllTasks()
				}
			}, win)
		}
	})

	list := widget.NewListWithData(
		vt.untypedList,
		func() fyne.CanvasObject {
			return component.NewTaskItem("", nil)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			task, _ := item.(binding.Untyped).Get()
			t := task.(m.Task)

			taskItem := obj.(*component.TaskItem)
			taskItem.SetTask(t)

			doDelete := func() {
				delErr := vt.service.DeleteTask(t.ID)
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
				if util.GetBoolPrefByKey(constant2.DevModePref) {
					doDelete()
				} else {
					dialog.ShowConfirm("Удаление задания", "Вы действительно хотите удалить задание?", func(b bool) {
						if b {
							doDelete()
						}
					}, win)
				}
			}
		})
	scroll := container.NewVScroll(list)
	controlBox := container.NewHBox(modeLbl, refreshBtn, addBtn, quickAddBtn, layout.NewSpacer(), deleteAllBtn)
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
