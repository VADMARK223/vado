package tasks

import (
	"time"
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
		ID: 0,
		//ID:        -1,
		Name:      util.Tpl("Fast %s task", storeMode),
		Completed: util.RndBool(),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(t model.Task) {
	now := time.Now()
	//if t.ID == -1 { // Новая задача
	if t.ID == 0 { // Новая задача
		t.CreatedAt = now
		t.UpdatedAt = now
	} else {
		t.UpdatedAt = now
	}
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
