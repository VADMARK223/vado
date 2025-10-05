package settings

import (
	"image/color"
	"os"
	"os/exec"
	"vado/internal/constant"
	"vado/internal/gui/common"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateView представление вкладки настроек приложения
//
// При локальной отладке сохраняет в /home/vadmark/.var/app/com.jetbrains.GoLand/config/fyne/io.vado
func CreateView() fyne.CanvasObject {
	prefs := fyne.CurrentApp().Preferences()

	autoStartServerHTTPCheck := widget.NewCheck("Стартовать HTTP-сервер при загрузке", func(checked bool) {
		prefs.SetBool(constant.AutoStartServerHTTP, checked)
	})
	autoStartServerHTTPCheck.Checked = prefs.Bool(constant.AutoStartServerHTTP)

	autoStartServerGRPCCheck := widget.NewCheck("Стартовать gRPC-сервер при загрузке", func(checked bool) {
		prefs.SetBool(constant.AutoStartServerGRPC, checked)
	})
	autoStartServerGRPCCheck.Checked = prefs.Bool(constant.AutoStartServerGRPC)

	textMode := widget.NewLabel("Тип хранилища:")
	toggleMode := NewToggleDefault()
	toggleMode.SetToggle(prefs.Bool(constant.StoreModePref))
	toggleMode.SetOffBtnText("DB")
	toggleMode.SetOnBtnText("JSON")
	toggleMode.OnChanged = func(v *bool) {
		switch {
		case v == nil:
			logger.L().Warn("Mode not set.")
		case *v:
			prefs.SetBool(constant.StoreModePref, true)
		default:
			prefs.SetBool(constant.StoreModePref, false)
		}
	}
	tipText := canvas.NewText("(после изменения перезапустите приложение)", color.White)
	tipText.TextStyle = fyne.TextStyle{Italic: true}
	restartBtn := common.NewBtn("Перезапуск", nil, func() {
		exe, _ := os.Executable()
		cmd := exec.Command(exe)
		err := cmd.Start()
		if err != nil {
			return
		}
		os.Exit(0)
	})

	boxMode := container.NewHBox(textMode, toggleMode, tipText, restartBtn)

	content := container.NewVBox(NewFastModeCheck(true), autoStartServerHTTPCheck, autoStartServerGRPCCheck, boxMode)
	return content
}
