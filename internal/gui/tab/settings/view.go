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

	autoStartServerHTTPCheck := widget.NewCheck("Стартовать HTTP-сервер при загрузке", func(checked bool) {
		prefs.SetBool(constant.AutoStartServerHTTP, checked)
	})
	autoStartServerHTTPCheck.Checked = prefs.Bool(constant.AutoStartServerHTTP)

	autoStartServerGRPCCheck := widget.NewCheck("Стартовать gRPC-сервер при загрузке", func(checked bool) {
		prefs.SetBool(constant.AutoStartServerGRPC, checked)
	})
	autoStartServerGRPCCheck.Checked = prefs.Bool(constant.AutoStartServerGRPC)

	content := container.NewVBox(devModeCheck, autoStartServerHTTPCheck, autoStartServerGRPCCheck)
	return content
}
