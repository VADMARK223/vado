package tasks

import (
	"vado/internal/model"
	"vado/pkg/util"
)

func (vt *ViewTasks) AddTaskQuick(isJSON bool) {
	storeMode := func() string {
		if isJSON {
			return "json"
		}
		return "db"
	}()
	newTask := model.Task{
		ID:        -1,
		Name:      util.Tpl("Fast %s task", storeMode),
		Completed: util.RndBool(),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(newTask model.Task) {
	_ = vt.service.CreateTask(newTask)
	_ = vt.updateUntypedList()
}

func (vt *ViewTasks) DeleteAllTasks() {
	vt.service.DeleteAllTasks()
	_ = vt.updateUntypedList()
}

func (vt *ViewTasks) GetTaskByID(id int) (*model.Task, error) {
	return vt.service.GetTaskByID(id)
}
