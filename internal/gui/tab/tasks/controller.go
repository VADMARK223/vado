package tasks

import (
	"vado/internal/model"
	"vado/pkg/util"
)

func (vt *ViewTasks) AddTaskQuick(isJSON bool) {
	id := util.RndIntn(10000)
	storeMode := func() string {
		if isJSON {
			return "json"
		}
		return "db"
	}()
	newTask := model.Task{
		ID:        id,
		Name:      util.Tpl("Fast %s task %d", storeMode, id),
		Completed: util.RndBool(),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(newTask model.Task) {
	_ = vt.service.CreateTask(newTask)
	_ = vt.reloadTasks()
}

func (vt *ViewTasks) DeleteAllTasks() {
	vt.service.DeleteAllTasks()
	_ = vt.reloadTasks()
}

func (vt *ViewTasks) GetTaskByID(id int) (model.Task, error) {
	return vt.service.GetTaskByID(id)
}
