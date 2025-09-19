package repository

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"vado/gui/tab/tasks/constant"
	"vado/model"
)

type TaskRepository struct {
	FilePath string
}

func (r *TaskRepository) LoadTasksList() (model.TaskList, error) {
	data, err := os.ReadFile(constant.TaskFilePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list model.TaskList
	err = json.Unmarshal(data, &list)
	return list, err
}

func (r *TaskRepository) SaveTasks(tasksList model.TaskList) error {
	err := os.MkdirAll(filepath.Dir(r.FilePath), 0755)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasksList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.FilePath, data, 0644)
}
