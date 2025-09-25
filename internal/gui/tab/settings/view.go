package settings

import (
	"vado/internal/constant"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateView представление вкладки настроек приложения
//
// При локальной отладке сохраняет в /home/vadmark/.var/app/com.jetbrains.GoLand/config/fyne/io.vado
func CreateView() fyne.CanvasObject {
	prefs := fyne.CurrentApp().Preferences()

	devModeCheck := widget.NewCheck("Режим разработчика", func(checked bool) {
		prefs.SetBool(constant.DevModePref, checked)
	})
	devModeCheck.Checked = prefs.Bool(constant.DevModePref)
	content := container.NewVBox(devModeCheck)
	return content
}
