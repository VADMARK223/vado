package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/model"
)

type TaskJSONRepo struct {
	filePath string
}

func (r *TaskJSONRepo) GetTask(id int) (model.Task, error) {
	//TODO implement me
	fmt.Println(id)
	panic("implement me for get by id: ")
}

func NewTaskJSONRepo(path string) *TaskJSONRepo {
	return &TaskJSONRepo{filePath: path}
}

func (r *TaskJSONRepo) FetchAll() (model.TaskList, error) {
	data, err := os.ReadFile(constant.TasksFilePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list model.TaskList
	err = json.Unmarshal(data, &list)
	return list, err
}

func (r *TaskJSONRepo) Save(t model.Task) error {
	list, _ := r.FetchAll()
	for _, task := range list.Tasks {
		if task.ID == t.ID {
			return errors.New("task already exists")
		}
	}
	list.Tasks = append(list.Tasks, t)

	return r.SaveTasks(list)
}

func (r *TaskJSONRepo) Remove(id int) error {
	list, _ := r.FetchAll()
	newTasks := make([]model.Task, 0)
	for _, t := range list.Tasks {
		if id != t.ID {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks
	return r.SaveTasks(list)
}

func (r *TaskJSONRepo) RemoveAll() error {
	list, _ := r.FetchAll()
	list.Tasks = []model.Task{}
	return r.SaveTasks(list)
}
