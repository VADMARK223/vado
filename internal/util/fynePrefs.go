package util

import (
	"fyne.io/fyne/v2"
)

func GetBoolPrefByKey(key string) bool {
	return fyne.CurrentApp().Preferences().Bool(key)
}
