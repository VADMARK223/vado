package util

import (
	"fmt"
	"math/rand"

	"fyne.io/fyne/v2"
)

func Tpl(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func GetBoolPrefByKey(key string) bool {
	return fyne.CurrentApp().Preferences().Bool(key)
}

func GetRandomBool() bool {
	return rand.Intn(2) == 1
}
