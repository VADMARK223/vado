package constant

import (
	"fmt"
	"os"
	"vado/pkg/logger"

	"github.com/joho/godotenv"
)

const TasksFilePath = "./data/tasks.json"

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		logger.L().Info("APP_ENV is empty.")
	}

	if env == "dev" {
		envFile := fmt.Sprintf("vado_app.env.%s", env)
		if err := godotenv.Load(envFile); err != nil {
			panic(fmt.Sprintf("Failed to load config %s: %v", envFile, err))
		}
	}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("Env variable %s not found", key))
}

// GetDSN Data Source Name (Имя источника данных)
func GetDSN() string {
	var host string
	host = getEnv("DB_HOST")
	port := getEnv("DB_PORT")
	user := getEnv("DB_USER")
	password := getEnv("DB_PASSWORD")
	dbname := getEnv("DB_NAME")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}
