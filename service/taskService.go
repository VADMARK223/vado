package service

import (
	"errors"
	"fmt"
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

func GetTaskById(id int) {
	fmt.Println("Get Task By Id")
}

func DeleteTask() {

}

func UpdateTask() {
	fmt.Println("Update Task")
}
