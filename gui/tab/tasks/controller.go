package tasks

import (
	"math/rand"
	"vado/model"
)

func (vt *ViewTasks) AddTaskFast(name string) {
	id := rand.Intn(10000)
	newTask := model.Task{Id: id, Name: name}
	err := vt.service.Create(newTask)
	if err != nil {
		return
	}

	vt.reloadTasks()
}
