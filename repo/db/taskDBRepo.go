package db

import (
	"database/sql"
	"fmt"
	"log"
	"vado/model"

	_ "github.com/lib/pq"
)

type TaskDBRepo struct {
	dataSourceName string
	db             *sql.DB
}

func NewTaskDBRepo(dataSourceName string) *TaskDBRepo {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Err open sql", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error ping connection: ", err)
	}

	fmt.Println("Successfully database connected!")

	return &TaskDBRepo{dataSourceName: dataSourceName, db: db}
}

func (r TaskDBRepo) FetchAll() (model.TaskList, error) {
	rows, err := r.db.Query("SELECT id, name, description, completed FROM tasks")
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
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Completed)
		if err != nil {
			log.Fatal("Task rows scan error:", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	var list model.TaskList
	list.Tasks = tasks

	return list, nil
}

func (r *TaskDBRepo) Save(task model.Task) error {
	query := `INSERT INTO tasks (id, name, description, completed) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, task.ID, task.Name, task.Description, task.Completed)
	if err != nil {
		return fmt.Errorf("error saving task: %w", err)
	}
	return nil
}

func (r TaskDBRepo) Remove(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting task with id %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}

	return nil
}

func (r TaskDBRepo) RemoveAll() error {
	// Выполняем удаление
	_, err := r.db.Exec("TRUNCATE TABLE tasks CASCADE")
	if err != nil {
		return fmt.Errorf("error delete all tasks: %w", err)
	}

	fmt.Println("All tasks removed successfully")
	return nil
}
