package database

import (
	"database/sql"
	"fmt"
	"log"
	"vado/model"

	_ "github.com/lib/pq" // sql: unknown driver "postgres" (forgotten import?)
)

func RunDatabase() {
	fmt.Println("Start database...")

	host := "127.0.0.1"
	port := 5432
	user := "vadmark"
	password := "5125341"
	dbname := "vadodb"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Подключение к базе
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")

	rows, err := db.Query("SELECT id, name, description FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.Id, &t.Name, &t.Description)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, t := range tasks {
		fmt.Printf("%+v\n", t)
	}
}
