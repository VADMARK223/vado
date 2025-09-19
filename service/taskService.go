package service

import (
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
	fmt.Println("Get All Tasks")
	return s.repo.LoadTasks()
}

func GetTaskById(id int) {
	fmt.Println("Get Task By Id")
}

func DeleteTask() {

}

func UpdateTask() {
	fmt.Println("Update Task")
}
