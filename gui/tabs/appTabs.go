package tabs

import (
	"vado/gui/common"

	"fyne.io/fyne/v2/container"
)

func CreateAppTabs() *container.AppTabs {
	settingsTabItem := common.CreateTabItem("Settings", CreateSettingsGui())
	settingsTabItem.Disabled()

	tabs := container.NewAppTabs(
		common.CreateTabItem("Http", CreateHttpTab()),
		common.CreateTabItem("Modules", CreateModulesTab()),
		common.CreateTabItem("Tasks", CreateTasksTab()),
		settingsTabItem,
	)
	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}
