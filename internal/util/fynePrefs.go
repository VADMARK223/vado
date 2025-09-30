package util

import (
	"vado/internal/constant"

	"fyne.io/fyne/v2"
)

// /home/vadmark/.var/app/com.jetbrains.GoLand/config/fyne/io.vado

func GetBoolPrefByKey(key string) bool {
	if fyne.CurrentApp() == nil {
		panic("Current app is nil!")
	}
	return fyne.CurrentApp().Preferences().Bool(key)
}

func IsFastMode() bool {
	if !IsDevMode() {
		return IsDevMode()
	}
	if fyne.CurrentApp() == nil {
		return false
	}

	return GetBoolPrefByKey(constant.FastModePref)
}

func AutoStartServerHTTP() bool {
	if fyne.CurrentApp() == nil {
		return false
	}

	return GetBoolPrefByKey(constant.AutoStartServerHTTP)
}

func AutoStartServerGRPC() bool {
	if fyne.CurrentApp() == nil {
		return false
	}

	return GetBoolPrefByKey(constant.AutoStartServerGRPC)
}

func OnFastModeChange(callback func(newValue bool)) {
	if fyne.CurrentApp() == nil {
		return
	}
	pref := fyne.CurrentApp().Preferences()
	pref.AddChangeListener(func() {
		newValue := GetBoolPrefByKey(constant.FastModePref)
		callback(newValue)
	})
}
