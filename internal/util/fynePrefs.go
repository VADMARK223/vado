package util

import (
	"vado/internal/constant"

	"fyne.io/fyne/v2"
)

func GetBoolPrefByKey(key string) bool {
	if fyne.CurrentApp() == nil {
		panic("Current app is nil!")
	}
	return fyne.CurrentApp().Preferences().Bool(key)
}

func IsDevMode() bool {
	if fyne.CurrentApp() == nil {
		return false
	}

	return GetBoolPrefByKey(constant.DevModePref)
}

func AutoStartServer() bool {
	if fyne.CurrentApp() == nil {
		return false
	}

	return GetBoolPrefByKey(constant.AutoStartServer)
}

func OnDevModeChange(callback func(newValue bool)) {
	if fyne.CurrentApp() == nil {
		return
	}
	pref := fyne.CurrentApp().Preferences()
	pref.AddChangeListener(func() {
		newValue := GetBoolPrefByKey(constant.DevModePref)
		callback(newValue)
	})
}
