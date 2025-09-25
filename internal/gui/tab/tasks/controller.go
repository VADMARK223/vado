package tasks

import (
	"math/rand"
	"vado/internal/model"
	"vado/pkg/util"
)

func (vt *ViewTasks) AddTaskQuick() {
	id := rand.Intn(10000)
	newTask := model.Task{
		ID:        id,
		Name:      util.Tpl("Fast task %d", id),
		Completed: util.GetRandomBool(),
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
