package grpc

import (
	"context"
	"vado/internal/domain/user"
	"vado/internal/pb/userpb"
)

// UserGRPCHandler реализует gRPC сервис для пользователей
type UserGRPCHandler struct {
	userpb.UnimplementedCreateUserServer
	service *user.Service
}

func NewUserGRPCHandler(service *user.Service) *UserGRPCHandler {
	return &UserGRPCHandler{service: service}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	domainUser := user.User{
		Username: req.Username,
		Password: req.Password,
	}

	err := h.service.CreateUser(domainUser)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:       int64(domainUser.ID),
			Username: domainUser.Username,
		},
	}, nil
}
