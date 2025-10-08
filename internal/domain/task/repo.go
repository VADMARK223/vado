package task

type Repo interface {
	FetchAll() (List, error)
	InsertUpdate(task Task) error
	Remove(id int) error
	GetTask(id int) (*Task, error)
	RemoveAll() error
}
