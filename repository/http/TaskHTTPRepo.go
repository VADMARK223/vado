package http

import "vado/model"

type TaskHTTPRepo struct {
	baseURL string
}

func NewTaskHTTPRepo(baseURL string) *TaskHTTPRepo {
	return &TaskHTTPRepo{baseURL: baseURL}
}

func (r TaskHTTPRepo) GetAllTasks() (model.TaskList, error) {
	//TODO implement me
	panic("implement me")
}

func (r TaskHTTPRepo) CreateTask(t model.Task) error {
	//TODO implement me
	panic("implement me")
}

func (r TaskHTTPRepo) DeleteTask(id int) error {
	//TODO implement me
	panic("implement me")
}

func (r TaskHTTPRepo) DeleteAllTasks() {
	//TODO implement me
	panic("implement me")
}
