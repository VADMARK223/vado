package model

// Task модель задания.
type Task struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Completed bool   `json:"completed"`
}

// TaskList модель списка заданий.
type TaskList struct {
	Tasks []Task `json:"tasks"`
}
