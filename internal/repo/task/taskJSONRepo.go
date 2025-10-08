package task

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
	"vado/internal/gui/tab/tasks/constant"
	"vado/internal/model"
	"vado/internal/util"
)

type TaskJSONRepo struct {
	filePath string
}

func (r *TaskJSONRepo) GetTask(id int) (*model.Task, error) {
	data, err := os.ReadFile(constant.TasksFilePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list model.TaskList
	err = json.Unmarshal(data, &list)

	for _, task := range list.Tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, err
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

func (r *TaskJSONRepo) InsertUpdate(t model.Task) error {
	list, _ := r.FetchAll()
	var now = time.Now()
	if t.ID == 0 { // новая задача
		t.ID = util.GenerateMaxID(list)
		t.CreatedAt = &now
		t.UpdatedAt = &now
		list.Tasks = append(list.Tasks, t)
	} else {
		updated := false
		for i, task := range list.Tasks {
			if task.ID == t.ID {
				t.UpdatedAt = &now
				list.Tasks[i] = t
				updated = true
				break
			}
		}

		if !updated {
			return errors.New("task not exists")
		}
	}

	return r.saveTasks(list)
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
	return r.saveTasks(list)
}

func (r *TaskJSONRepo) RemoveAll() error {
	list, _ := r.FetchAll()
	list.Tasks = []model.Task{}
	return r.saveTasks(list)
}

func (r *TaskJSONRepo) saveTasks(tasksList model.TaskList) error {
	err := os.MkdirAll(filepath.Dir(r.filePath), 0755)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasksList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}
