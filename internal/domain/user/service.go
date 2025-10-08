package user

type Service struct {
	repo *DBRepo
}

func NewUserService(repo *DBRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(user User) error {
	return s.repo.CreateUser(user)
}
