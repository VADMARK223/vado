package tab

import (
	"vado/gui/common"
	"vado/gui/tab/http"
	"vado/gui/tab/todo"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const defaultTabIndex = 1

func CreateAppTabs(win fyne.Window) *container.AppTabs {
	factories := map[*container.TabItem]func() fyne.CanvasObject{}

	tabs := container.NewAppTabs(
		common.CreateLazyTabItem("Http", http.CreateHttpTab, factories),
		common.CreateLazyTabItem("TODO", func() fyne.CanvasObject {
			return todo.CreateTODOTab(win)
		}, factories),
		common.CreateLazyTabItem("Модули", CreateModulesTab, factories),
		common.CreateLazyTabItem("Настройки", CreateSettingsTab, factories),
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
