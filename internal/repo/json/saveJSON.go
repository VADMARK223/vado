package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"vado/internal/model"
)

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
