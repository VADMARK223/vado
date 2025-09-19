package service

import (
	"errors"
	"vado/model"
	"vado/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() (model.TaskList, error) {
	return s.repo.LoadTasksList()
}

func (s *TaskService) Create(t model.Task) error {
	list, _ := s.repo.LoadTasksList()
	for _, task := range list.Tasks {
		if task.Id == t.Id {
			return errors.New("task already exists")
		}
	}
	list.Tasks = append(list.Tasks, t)

	return s.repo.SaveTasks(list)
}

// Delete удаляет задание по его идентификатору
// TODO: добавить проверку на отсутствие задания с таким идентификатором
func (s *TaskService) Delete(id int) error {
	list, _ := s.repo.LoadTasksList()
	newTasks := make([]model.Task, 0)
	for _, t := range list.Tasks {
		if id != t.Id {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks
	return s.repo.SaveTasks(list)
}

func (s *TaskService) DeleteAllTasks() {
	list, _ := s.repo.LoadTasksList()
	list.Tasks = []model.Task{}
	err := s.repo.SaveTasks(list)
	if err != nil {
		panic(err)
	}
}
