package service

import (
	"errors"
	"vado/internal/model"
	"vado/internal/repo"
)

var ErrInvalidTask = errors.New("invalid task")

type ITaskService interface {
	GetAllTasks() (model.TaskList, error)
	CreateTask(t model.Task) error
	GetTaskByID(id int) (*model.Task, error)
	DeleteTask(id int) error
	DeleteAllTasks()
}

// TaskService конкретная реализация бизнес-логики
type TaskService struct {
	repo repo.TaskRepo
}

func NewTaskService(repo repo.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() (model.TaskList, error) {
	return s.repo.FetchAll()
}

func (s *TaskService) CreateTask(task model.Task) error {
	return s.repo.InsertUpdate(task)
}

func (s *TaskService) GetTaskByID(id int) (*model.Task, error) {
	return s.repo.GetTask(id)
}

func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Remove(id)
}

func (s *TaskService) DeleteAllTasks() {
	err := s.repo.RemoveAll()
	if err != nil {
		return
	}
}
