package tasks

import (
	"fmt"
	"image/color"
	"vado/internal/gui/common"
	c "vado/internal/gui/tab/tasks/component"
	"vado/internal/gui/tab/tasks/component/grpc"
	"vado/internal/gui/tab/tasks/component/http"
	m "vado/internal/model"
	"vado/internal/service"
	util2 "vado/internal/util"
	"vado/pkg/logger"
	"vado/pkg/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

type ViewTasks struct {
	service     service.ITaskService
	untypedList binding.UntypedList
}

func NewTasksView(win fyne.Window, s service.ITaskService) fyne.CanvasObject {
	vt := &ViewTasks{service: s, untypedList: binding.NewUntypedList()}
	err := vt.updateUntypedList()
	if err != nil {
		return nil
	}

	isJSON := util2.IsJSONMode()

	modeLbl := widget.NewRichTextFromMarkdown(func() string {
		var mode string
		if isJSON {
			mode = "JSON"
		} else {
			mode = "DB"
		}

		return fmt.Sprintf("Источник: **%s**", mode)
	}())

	refreshBtn := common.NewBtn("", theme.ViewRefreshIcon(), func() {
		_ = vt.updateUntypedList()
	})
	addBtn := common.NewBtn("", theme.ContentAddIcon(), func() {
		showTaskDialog(win, vt, nil)
	})
	quickAddBtn := common.NewBtn("Быстро", theme.ContentAddIcon(), func() {
		vt.AddTaskQuick(isJSON)
	})
	updateQuickAddBtnVisibility := func() {
		if util2.IsFastMode() {
			quickAddBtn.Show()
		} else {
			quickAddBtn.Hide()
		}
	}

	updateQuickAddBtnVisibility()

	util2.OnFastModeChange(func(newValue bool) {
		updateQuickAddBtnVisibility()
	})

	deleteAllBtn := common.NewBtn("Удалить все", theme.DeleteIcon(), func() {
		if util2.IsFastMode() {
			vt.DeleteAllTasks()
		} else {
			dialog.ShowConfirm("Удаление всех заданий", "Вы действительно хотите удалить все задания?", func(b bool) {
				if b {
					vt.DeleteAllTasks()
				}
			}, win)
		}
	})

	list := widget.NewListWithData(
		vt.untypedList,
		func() fyne.CanvasObject {
			return c.NewTaskItem("", nil)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			task, _ := item.(binding.Untyped).Get()
			t := task.(m.Task)

			taskItem := obj.(*c.TaskItem)
			taskItem.SetTask(t)

			doDelete := func() {
				delErr := vt.service.DeleteTask(t.ID)
				if delErr != nil {
					logger.L().Error("Failed to delete task", zap.Error(delErr))
					dialog.ShowError(delErr, win)
					return
				}
				_ = vt.updateUntypedList()
				if err != nil {
					logger.L().Error("Failed to update list", zap.Error(err))
					return
				}
			}

			taskItem.OnDoubleTap = func() {
				requestedTask, err := vt.GetTaskByID(t.ID)
				if err != nil {
					logger.L().Info(fmt.Sprintf("Get task %d", t.ID), zap.String("Error: ", err.Error()))
				}
				showTaskDialog(win, vt, requestedTask)
			}

			taskItem.OnDelete = func() {
				if util2.IsFastMode() {
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
	controlBox := container.NewHBox(modeLbl, refreshBtn, addBtn, layout.NewSpacer(), quickAddBtn, deleteAllBtn)
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
	topBox := container.NewVBox(http.NewControlBoxHTTP(vt.service), grpc.NewControlBoxGRPC(vt.service, win), controlBox, title)
	content := container.NewBorder(topBox, nil, nil, nil, scroll)
	return content

}

func showTaskDialog(win fyne.Window, vt *ViewTasks, t *m.Task) {
	c.ShowTaskDialog(win, t, func(task m.Task) {
		vt.AddTask(task)

		err := vt.updateUntypedList()
		if err != nil {
			logger.L().Error("Error reload tasks.")
		}

		err = vt.updateUntypedList()
		if err != nil {
			logger.L().Error("Error reset list.")
		}
	})
}

// resetUntypedList метод пересоздает весь список, потому что изменения деталей задач не отображаются, после редактирования
func (vt *ViewTasks) updateUntypedList() error {
	cur, _ := vt.service.GetAllTasks()
	_ = vt.untypedList.Set([]any{})

	items := make([]any, len(cur.Tasks))
	for i, t := range cur.Tasks {
		items[i] = t
	}

	return vt.untypedList.Set(items)
}
