package util

func GetModeValue() string {
	return "dev" //os.Getenv("APP_ENV")
}

func IsDevMode() bool {
	env := GetModeValue()
	if env == "" {
		return true
	}

	return env == "dev"
}
