package tasks

import (
	"vado/internal/model"
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
	newTask := model.Task{
		ID:        0,
		Name:      util.Tpl("Fast %s task", storeMode),
		Completed: util.RndBool(),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(t model.Task) {
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

func (vt *ViewTasks) GetTaskByID(id int) (*model.Task, error) {
	return vt.service.GetTaskByID(id)
}
