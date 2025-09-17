package tabs

import (
	"vado/gui/common"
	"vado/gui/tabs/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateAppTabs() *container.AppTabs {
	factories := map[*container.TabItem]func() fyne.CanvasObject{}

	tabs := container.NewAppTabs(
		common.CreateLazyTabItem("Http", http.CreateHttpTab, factories),
		common.CreateLazyTabItem("Tasks", CreateTasksTab, factories),
		common.CreateLazyTabItem("Modules", CreateModulesTab, factories),
		common.CreateLazyTabItem("Settings", CreateSettingsTab, factories),
	)
	tabs.SelectIndex(1)
	tabs.SetTabLocation(container.TabLocationTop)

	loaded := map[*container.TabItem]bool{}
	tabs.OnSelected = func(item *container.TabItem) {
		if item == nil || loaded[item] {
			return
		}

		if factory, ok := factories[item]; ok {
			item.Content = factory()
		}

		loaded[item] = true
		delete(factories, item)
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
