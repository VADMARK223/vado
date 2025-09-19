package tasks

import (
	"vado/model"
)

func (vt *ViewTasks) AddTask(newTask model.Task) {
	_ = vt.service.Create(newTask)
	_ = vt.reloadTasks()
}

func (vt *ViewTasks) DeleteAllTasks() {
	vt.service.DeleteAllTasks()
	_ = vt.reloadTasks()
}
