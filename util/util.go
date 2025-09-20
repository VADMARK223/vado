package util

import (
	"fmt"

	"fyne.io/fyne/v2"
)

func Tpl(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func GetBoolPrefByKey(key string) bool {
	return fyne.CurrentApp().Preferences().Bool(key)

}
