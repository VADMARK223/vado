package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"vado/internal/model"
	"vado/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type TaskDBRepo struct {
	dataSourceName string
	db             *sql.DB
}

func NewTaskDBRepo(dsn string) *TaskDBRepo {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Err open sql", err)
	}

	logger.L().Info("Try connect to database...", zap.String("dsn", dsn))
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error ping connection: %v", err))
	}

	logger.L().Info(fmt.Sprintf("Successfully database connected! (%s)", dsn))

	return &TaskDBRepo{dataSourceName: dsn, db: db}
}

func (t *TaskDBRepo) FetchAll() (model.TaskList, error) {
	rows, err := t.db.Query(`SELECT id, name, description, completed FROM tasks ORDER BY created_at`)
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

func (t *TaskDBRepo) Save(task model.Task) error {
	if task.ID == -1 {
		// Новая задача — вставляем
		query := `
			INSERT INTO tasks (name, description, completed)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		return t.db.QueryRow(query, task.Name, task.Description, task.Completed).Scan(&task.ID)
	} else {
		// Существующая задача — обновляем
		query := `
			UPDATE tasks
			SET name = $1, description = $2, completed = $3
			WHERE id = $4
		`
		_, err := t.db.Exec(query, task.Name, task.Description, task.Completed, task.ID)
		return err
	}
}

func (t *TaskDBRepo) GetTask(id int) (*model.Task, error) {
	query := `SELECT id, name, description, completed FROM tasks WHERE id = $1`
	var task model.Task
	err := t.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Completed,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// если задачи нет → вернем понятную ошибку
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &task, nil
}

func (t *TaskDBRepo) Remove(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := t.db.Exec(query, id)
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

func (t *TaskDBRepo) RemoveAll() error {
	_, err := t.db.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("error delete all tasks: %w", err)
	}

	logger.L().Debug("Successfully deleted all tasks.")
	return nil
}
