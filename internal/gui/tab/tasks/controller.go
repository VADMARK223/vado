package tasks

import (
	"vado/internal/model"
	"vado/pkg/util"
)

func (vt *ViewTasks) AddTaskQuick() {
	id := util.RndIntn(10000)
	newTask := model.Task{
		ID:        id,
		Name:      util.Tpl("Fast task %d", id),
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
