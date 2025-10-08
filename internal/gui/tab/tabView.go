package tab

import (
	task2 "vado/internal/domain/task"
	"vado/internal/gui/common"
	"vado/internal/gui/tab/admin"
	"vado/internal/gui/tab/heavy"
	"vado/internal/gui/tab/lesson"
	"vado/internal/gui/tab/settings"
	"vado/internal/gui/tab/tasks"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/server/context"
	"vado/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const defaultTabIndex = 2

func NewTabsView(appCtx *context.AppContext, win fyne.Window) *container.AppTabs {
	factories := map[*container.TabItem]func() fyne.CanvasObject{}

	tabs := container.NewAppTabs(
		common.CreateLazyTabItem("Задания", func() fyne.CanvasObject {
			var r task2.TaskRepo
			if util.IsJSONMode() {
				r = task2.NewTaskJSONRepo(constant.TasksFilePath)
			} else {
				r = task2.NewTaskDBRepo(appCtx.DB)
			}
			appCtx.HttpContext.TaskService.Repo = r
			//s := task2.NewTaskService(r)
			return tasks.NewTasksView(appCtx, win)
		}, factories),
		common.CreateLazyTabItem("Настройки", settings.CreateView, factories),
		common.CreateLazyTabItem("Админка", func() fyne.CanvasObject {
			return admin.NewAdminView(appCtx, win)
		}, factories),
		common.CreateLazyTabItem("Уроки", lesson.CreateView, factories),
		common.CreateLazyTabItem("Тяжелая вкладка", heavy.NewHeavyView, factories),
	)
	tabs.SelectIndex(defaultTabIndex)
	tabs.SetTabLocation(container.TabLocationTop)

	loaded := map[*container.TabItem]bool{}
	tabs.OnSelected = func(item *container.TabItem) {
		if item == nil || loaded[item] {
			return
		}

		if factory, ok := factories[item]; ok {
			go func() {
				// Создаем контент в фоне
				content := factory()
				// Обновляем UI безопасно
				fyne.Do(func() {
					item.Content = content

					loaded[item] = true
					delete(factories, item)
				})
			}()
		}
	}

	// На всякий случай — если SelectIndex не вызвал OnSelected (в некоторых версиях/сценариях),
	// Подгружаем контент для выбранной вкладки вручную:
	if sel := tabs.Selected(); sel != nil && !loaded[sel] {
		if f, ok := factories[sel]; ok {
			sel.Content = f()
			if sel.Content != nil {
				sel.Content.Refresh()
			}
			loaded[sel] = true
			delete(factories, sel)
		}
	}

	return tabs
}
