package util

import "fmt"

func Tpl(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
