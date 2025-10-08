package tasks

import (
	"vado/internal/domain/task"
	"vado/pkg/logger"
	"vado/pkg/util"

	"go.uber.org/zap"
)

func (vt *ViewTasks) AddTaskQuick(isJSON bool) {
	storeMode := func() string {
		if isJSON {
			return "json"
		}
		return "db"
	}()
	newTask := task.Task{
		ID:        0,
		Name:      util.Tpl("Fast %s task", storeMode),
		Completed: util.RndBool(),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(t task.Task) {
	err := vt.service.CreateTask(t)
	if err != nil {
		logger.L().Error("Error insert/update task", zap.Error(err))
	}
	_ = vt.updateUntypedList()
}

func (vt *ViewTasks) DeleteAllTasks() {
	vt.service.DeleteAllTasks()
	_ = vt.updateUntypedList()
}

func (vt *ViewTasks) GetTaskByID(id int) (*task.Task, error) {
	return vt.service.GetTaskByID(id)
}
