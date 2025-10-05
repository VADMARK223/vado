package util

import "vado/internal/model"

func GenerateMaxID(list model.TaskList) int {
	maxID := 0
	for _, task := range list.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}
