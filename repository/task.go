package repository

import (
	"encoding/json"
	"log"
	"os"
	"vado/gui/tab/tasks/constant"
	"vado/model"
)

type TaskRepository struct {
	FilePath string
}

func (r *TaskRepository) LoadTasks() (model.TaskList, error) {
	data, err := os.ReadFile(constant.TaskFilePath)

	if err != nil {
		log.Fatal("Error open file:", err)
	}

	var list model.TaskList
	err = json.Unmarshal(data, &list)
	return list, err
}
