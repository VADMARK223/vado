package constant

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const TasksFilePath = "./data/tasks.json"

func init() {
	if os.Getenv("APP_ENV") == "development" {
		fmt.Println("Load .env file")
		err := godotenv.Load("vado_db.env")
		if err != nil {
			panic("Error loading vado.env file!!!")
		}
	} else {
		fmt.Println("Not load .env file")
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
	if os.Getenv("APP_ENV") == "development" {
		host = "127.0.0.1"
	} else {
		host = getEnv("DB_HOST")
	}
	port := getEnv("DB_PORT")
	user := getEnv("DB_USER")
	password := getEnv("DB_PASSWORD")
	dbname := getEnv("DB_NAME")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}
