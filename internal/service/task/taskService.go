package task

import (
	"vado/internal/model"
	"vado/internal/repo"
)

type ITaskService interface {
	GetAllTasks() (model.TaskList, error)
	CreateTask(t model.Task) error
	GetTaskByID(id int) (*model.Task, error)
	DeleteTask(id int) error
	DeleteAllTasks()
}

type Service struct {
	Repo repo.TaskRepo
}

func NewTaskService(repo repo.TaskRepo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllTasks() (model.TaskList, error) {
	return s.Repo.FetchAll()
}

func (s *Service) CreateTask(task model.Task) error {
	return s.Repo.InsertUpdate(task)
}

func (s *Service) GetTaskByID(id int) (*model.Task, error) {
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
