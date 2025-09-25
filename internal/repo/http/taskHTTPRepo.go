package http

import (
	"encoding/json"
	"io"
	"net/http"
	"vado/internal/model"
)

type TaskHTTPRepo struct {
	baseURL string
	client  *http.Client
}

func NewTaskHTTPRepo(baseURL string) *TaskHTTPRepo {
	return &TaskHTTPRepo{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (t TaskHTTPRepo) FetchAll() (model.TaskList, error) {
	resp, err := t.client.Get(t.baseURL + "/tasks")
	if err != nil {
		return model.TaskList{}, err
	}
	defer resp.Body.Close()

	var list model.TaskList
	var tasks []model.Task
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &tasks); err != nil {
		return model.TaskList{}, err
	}
	list.Tasks = tasks
	return list, nil
}

func (t TaskHTTPRepo) Save(task model.Task) error {
	//TODO implement me
	panic("implement me")
}

func (t TaskHTTPRepo) Remove(id int) error {
	//TODO implement me
	panic("implement me")
}

func (t TaskHTTPRepo) RemoveAll() {
	//TODO implement me
	panic("implement me")
}
