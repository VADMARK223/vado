package tab

import (
	"vado/gui/common"
	"vado/gui/tab/http"
	"vado/gui/tab/lesson"
	"vado/gui/tab/settings"
	"vado/gui/tab/tasks"
	"vado/gui/tab/tasks/constant"
	"vado/repository"
	"vado/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const defaultTabIndex = 1

func CreateAppTabs(win fyne.Window) *container.AppTabs {
	factories := map[*container.TabItem]func() fyne.CanvasObject{}
	r := &repository.TaskRepository{FilePath: constant.TaskFilePath}
	s := service.NewTaskService(r)

	tabs := container.NewAppTabs(
		common.CreateLazyTabItem("Http", http.CreateView, factories),
		common.CreateLazyTabItem("Задания", func() fyne.CanvasObject {
			return tasks.CreateView(win, s)
		}, factories),
		common.CreateLazyTabItem("Уроки", lesson.CreateView, factories),
		common.CreateLazyTabItem("Настройки", settings.CreateView, factories),
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
