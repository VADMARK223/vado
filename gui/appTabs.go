package gui

import (
	"vado/gui/common"
	"vado/gui/modules"
	"vado/gui/settings"
	"vado/gui/tasks"

	"fyne.io/fyne/v2/container"
)

func CreateAppTabs() *container.AppTabs {
	settingsTabItem := common.CreateTabItem("Settings", settings.CreateSettingsGui())
	settingsTabItem.Disabled()

	tabs := container.NewAppTabs(
		common.CreateTabItem("Modules", modules.CreateModulesGui()),
		common.CreateTabItem("Tasks", tasks.CreateTasksGui()),
		settingsTabItem,
	)
	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}
