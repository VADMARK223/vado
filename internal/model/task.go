package model

// Task модель задачи.
// @Description Структура задачи.
type Task struct {
	ID          int    `json:"id" example:"1" format:"int64"`        // Уникальный идентификатор
	Name        string `json:"name" example:"Купить молоко"`         // Название задачи
	Description string `json:"description" example:"Купить 2 литра"` // Описание задачи
	Completed   bool   `json:"completed" example:"false"`            // Флаг выполнения
}

// TaskList модель списка задач.
type TaskList struct {
	Tasks []Task `json:"tasks"`
}
