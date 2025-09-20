package tasks

import (
	"math/rand"
	"vado/model"
	"vado/util"
)

func (vt *ViewTasks) AddTaskQuick() {
	id := rand.Intn(10000)
	newTask := model.Task{
		Id:        id,
		Name:      util.Tpl("Fast task %d", id),
		Completed: rand.Intn(2) == 1,
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
