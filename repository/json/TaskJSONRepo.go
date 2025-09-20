package json

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"vado/gui/tab/tasks/constant"
	"vado/model"
)

type TaskJSONRepo struct {
	filePath string
}

func NewTaskJSONRepo(path string) *TaskJSONRepo {
	return &TaskJSONRepo{filePath: path}
}

func (r *TaskJSONRepo) GetAllTasks() (model.TaskList, error) {
	data, err := os.ReadFile(constant.TasksFilePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list model.TaskList
	err = json.Unmarshal(data, &list)
	return list, err
}

func (r *TaskJSONRepo) CreateTask(t model.Task) error {
	list, _ := r.GetAllTasks()
	for _, task := range list.Tasks {
		if task.Id == t.Id {
			return errors.New("task already exists")
		}
	}
	list.Tasks = append(list.Tasks, t)

	return r.SaveTasks(list)
}

func (r *TaskJSONRepo) DeleteTask(id int) error {
	list, _ := r.GetAllTasks()
	newTasks := make([]model.Task, 0)
	for _, t := range list.Tasks {
		if id != t.Id {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks
	return r.SaveTasks(list)
}

func (r *TaskJSONRepo) DeleteAllTasks() {
	list, _ := r.GetAllTasks()
	list.Tasks = []model.Task{}
	err := r.SaveTasks(list)
	if err != nil {
		panic(err)
	}
}
