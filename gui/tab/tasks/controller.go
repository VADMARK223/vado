package tasks

import (
	"fmt"
	"log"
	"math/rand"
	"vado/model"
)

// GUI вызывает сервис: add, delete, refresh

func (vt *ViewTasks) AddTaskFast(name string) {
	fmt.Println("Add Task Fast")
	id := rand.Intn(10000)
	err := vt.service.Create(model.Task{Id: id, Name: name})
	if err != nil {
		log.Fatal("ZZZZZ")
		return
	}

	c, _ := vt.service.GetAllTasks()
	log.Println(len(c.Tasks))
	vt.reloadTasks()
	log.Println(">>>>", len(vt.tasks))
	vt.list.Refresh()
}
