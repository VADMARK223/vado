package settings

import (
	"fmt"
	"strings"
	"vado/internal/constant"
	"vado/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewFastModeCheck(inSettings bool) *widget.Check {
	prefs := fyne.CurrentApp().Preferences()

	devModeCheck := widget.NewCheck(func() string {
		if inSettings {
			var text = "Режим быстрой отладки GUI"
			text += fmt.Sprintf(" (Недоступно в %s)", strings.ToUpper(util.GetModeValue()))
			return text
		} else {
			return "Быстрая отладка"
		}

	}(), func(checked bool) {
		prefs.SetBool(constant.FastModePref, checked)
	})
	devModeCheck.Disable()

	util.OnFastModeChange(func(newValue bool) {
		devModeCheck.SetChecked(newValue)
	})

	devModeCheck.Checked = util.IsFastMode()
	return devModeCheck
}
