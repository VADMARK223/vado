package task

import (
	"vado/internal/domain/task"
)

type ITaskService interface {
	GetAllTasks() (task.TaskList, error)
	CreateTask(t task.Task) error
	GetTaskByID(id int) (*task.Task, error)
	DeleteTask(id int) error
	DeleteAllTasks()
}

type Service struct {
	Repo task.TaskRepo
}

func NewTaskService(repo task.TaskRepo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllTasks() (task.TaskList, error) {
	return s.Repo.FetchAll()
}

func (s *Service) CreateTask(task task.Task) error {
	return s.Repo.InsertUpdate(task)
}

func (s *Service) GetTaskByID(id int) (*task.Task, error) {
	return s.Repo.GetTask(id)
}

func (s *Service) DeleteTask(id int) error {
	return s.Repo.Remove(id)
}

func (s *Service) DeleteAllTasks() {
	err := s.Repo.RemoveAll()
	if err != nil {
		return
	}
}
