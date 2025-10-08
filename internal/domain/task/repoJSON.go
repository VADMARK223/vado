package task

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

type JSONRepo struct {
	filePath string
}

func (r *JSONRepo) GetTask(id int) (*Task, error) {
	data, err := os.ReadFile(r.filePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list List
	err = json.Unmarshal(data, &list)

	for _, task := range list.Tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, err
}

func NewTaskJSONRepo(path string) *JSONRepo {
	return &JSONRepo{filePath: path}
}

func (r *JSONRepo) FetchAll() (List, error) {
	data, err := os.ReadFile(r.filePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list List
	err = json.Unmarshal(data, &list)
	return list, err
}

func (r *JSONRepo) InsertUpdate(t Task) error {
	list, _ := r.FetchAll()
	var now = time.Now()
	if t.ID == 0 { // новая задача
		t.ID = generateMaxID(list)
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

func (r *JSONRepo) Remove(id int) error {
	list, _ := r.FetchAll()
	newTasks := make([]Task, 0)
	for _, t := range list.Tasks {
		if id != t.ID {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks
	return r.saveTasks(list)
}

func (r *JSONRepo) RemoveAll() error {
	list, _ := r.FetchAll()
	list.Tasks = []Task{}
	return r.saveTasks(list)
}

func (r *JSONRepo) saveTasks(tasksList List) error {
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

func generateMaxID(list List) int {
	maxID := 0
	for _, task := range list.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}
