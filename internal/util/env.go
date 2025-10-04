package util

import (
	"os"
)

func GetModeValue() string {
	return os.Getenv("APP_ENV")
}

func IsDevMode() bool {
	env := GetModeValue()
	if env == "" {
		panic("APP_ENV is empty.")
	}

	return env == "dev"
}
