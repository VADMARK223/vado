package user

import (
	"vado/internal/model"
	"vado/internal/repo/user"
)

type Service struct {
	repo *user.UserDBRepo
}

func NewUserService(repo *user.UserDBRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(user model.User) error {
	return s.repo.CreateUser(user)
}
