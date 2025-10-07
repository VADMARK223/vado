package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"vado/internal/model"
	"vado/pkg/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type TaskDBRepo struct {
	db *sql.DB
}

func NewTaskDBRepo(db *sql.DB) *TaskDBRepo {
	return &TaskDBRepo{db: db}
}

func (t *TaskDBRepo) FetchAll() (model.TaskList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := t.db.QueryContext(ctx, `SELECT id, name, description, completed, created_at, updated_at FROM tasks ORDER BY created_at`)
	if err != nil {
		logger.L().Error("Error query tasks:", zap.Error(err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.L().Error("Error query tasks:", zap.Error(err))
		}
	}(rows)

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			logger.L().Error("Task rows scan error:", zap.Error(err))
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		logger.L().Error("Task rows error:", zap.Error(err))
	}

	var list model.TaskList
	list.Tasks = tasks

	return list, nil
}

func (t *TaskDBRepo) InsertUpdate(task model.Task) error {
	if task.ID == 0 {
		// Новая задача — вставляем
		query := `
			INSERT INTO tasks (name, description, completed, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		// TODO 1 это id VADMARK
		return t.db.QueryRow(query, task.Name, task.Description, task.Completed, 1).Scan(&task.ID)
	} else {
		// Существующая задача — обновляем
		//*task.UpdatedAt = time.Now()
		query := `
			UPDATE tasks
			SET name = $1, description = $2, completed = $3, updated_at = $4
			WHERE id = $5
		`
		_, err := t.db.Exec(query, task.Name, task.Description, task.Completed, time.Now(), task.ID)
		return err
	}
}

func (t *TaskDBRepo) GetTask(id int) (*model.Task, error) {
	query := `SELECT id, name, description, completed, created_at, updated_at FROM tasks WHERE id = $1`
	var task model.Task
	err := t.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
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
