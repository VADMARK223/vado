package task

type ITaskService interface {
	GetAllTasks() (TaskList, error)
	CreateTask(t Task) error
	GetTaskByID(id int) (*Task, error)
	DeleteTask(id int) error
	DeleteAllTasks()
}

type Service struct {
	Repo TaskRepo
}

func NewTaskService(repo TaskRepo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllTasks() (TaskList, error) {
	return s.Repo.FetchAll()
}

func (s *Service) CreateTask(task Task) error {
	return s.Repo.InsertUpdate(task)
}

func (s *Service) GetTaskByID(id int) (*Task, error) {
	return s.Repo.GetTask(id)
}

func (s *Service) DeleteTask(id int) error {
	return s.Repo.Remove(id)
}

func (s *Service) DeleteAllTasks() {
	err := s.Repo.RemoveAll()
	if err != nil {
		return
	}
}
