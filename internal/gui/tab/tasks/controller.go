package tasks

import (
	"vado/internal/model"
	util2 "vado/internal/util"
	"vado/pkg/util"
)

func (vt *ViewTasks) AddTaskQuick(isJSON bool) {
	tasks, err := vt.service.GetAllTasks()
	if err != nil {
		return
	}
	id := util2.GenerateMaxID(tasks)
	storeMode := func() string {
		if isJSON {
			return "json"
		}
		return "db"
	}()
	newTask := model.Task{
		ID:        id,
		Name:      util.Tpl("Fast %s task", storeMode),
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

func (vt *ViewTasks) GetTaskByID(id int) (*model.Task, error) {
	return vt.service.GetTaskByID(id)
}
