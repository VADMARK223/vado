package model

// Task модель задания.
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// TaskList модель списка заданий.
type TaskList struct {
	Tasks []Task `json:"tasks"`
}
