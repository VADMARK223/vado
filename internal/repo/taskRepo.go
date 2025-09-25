package repo

import (
	"vado/internal/model"
)

type TaskRepo interface {
	FetchAll() (model.TaskList, error)
	Save(task model.Task) error
	Remove(id int) error
	RemoveAll() error
}
