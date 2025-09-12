package model

type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func NewTask(id int, name string, done bool) *Task {
	return &Task{id, name, done}
}
