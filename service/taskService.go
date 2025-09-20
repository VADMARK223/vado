package service

import (
	"vado/model"
)

type ITaskService interface {
	GetAllTasks() (model.TaskList, error)
	CreateTask(t model.Task) error
	DeleteTask(id int) error
	DeleteAllTasks()
}
