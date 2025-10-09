package server

import (
	"database/sql"
	"fmt"
	"log"
	"vado/internal/gui/tab/tasks/constant"
	"vado/pkg/logger"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dsn := constant.GetDSN()
	fmt.Printf("Try connect to database: %s\n", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Err open sql", err)
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error ping connection: %v", err))
	}

	logger.L().Info(fmt.Sprintf("Successfully database connected! (%s)", dsn))
	return db
}
