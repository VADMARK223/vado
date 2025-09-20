package tasks

import (
	"math/rand"
	"vado/model"
	"vado/util"
)

func (vt *ViewTasks) AddTaskQuick() {
	id := rand.Intn(10000)
	newTask := model.Task{
		Id:   id,
		Name: util.Tpl("Fast task %d", id),
	}
	vt.AddTask(newTask)
}

func (vt *ViewTasks) AddTask(newTask model.Task) {
	_ = vt.service.Create(newTask)
	_ = vt.reloadTasks()
}

func (vt *ViewTasks) DeleteAllTasks() {
	vt.service.DeleteAllTasks()
	_ = vt.reloadTasks()
}
