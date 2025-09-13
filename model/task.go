package model

type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}
