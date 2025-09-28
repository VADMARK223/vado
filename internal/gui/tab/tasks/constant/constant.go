package constant

import (
	"fmt"
	"os"
)

const TasksFilePath = "./data/tasks.json"

//var GetDSN = util.Tpl(
//	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
//	"127.0.0.1", 5432, "vadmark", "5125341", "vadodb",
//)

const (
	Host     = "127.0.0.1"
	Port     = 5432
	User     = "vadmark"
	Password = "5125341"
	DBName   = "vadodb"
)

func GetDSN() string {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "vadmark")
	password := getEnv("DB_PASSWORD", "5125341")
	dbname := getEnv("DB_NAME", "vadodb")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
