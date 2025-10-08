package task

type TaskRepo interface {
	FetchAll() (TaskList, error)
	InsertUpdate(task Task) error
	Remove(id int) error
	GetTask(id int) (*Task, error)
	RemoveAll() error
}
