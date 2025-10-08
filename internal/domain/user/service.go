package user

type Service struct {
	repo *UserDBRepo
}

func NewUserService(repo *UserDBRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(user User) error {
	return s.repo.CreateUser(user)
}
