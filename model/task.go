package model

type Task struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Completed bool   `json:"completed"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}
