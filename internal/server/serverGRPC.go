package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"vado/internal/domain/task"
	taskGRPC "vado/internal/domain/task/transport/grpc"
	"vado/internal/domain/user"
	userGRPC "vado/internal/domain/user/transport/grpc"
	"vado/internal/pb/taskpb"
	"vado/internal/pb/userpb"
	"vado/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Server    *grpc.Server
	port      string
	isRunning atomic.Bool

	userService *user.Service
	taskService task.ITaskService
}

func NewServerGRPC(userService *user.Service, taskService task.ITaskService, port string) *GRPCServer {
	return &GRPCServer{
		port:        port,
		userService: userService,
		taskService: taskService,
	}
}

func (s *GRPCServer) createServerInstance() {
	grpcServer := grpc.NewServer()

	// Регистрируем User Service
	userHandler := userGRPC.NewUserGRPCHandler(s.userService)
	userpb.RegisterCreateUserServer(grpcServer, userHandler)

	// Регистрируем Task Service
	taskHandler := taskGRPC.NewTaskGRPCHandler(s.taskService)
	taskpb.RegisterTaskServiceServer(grpcServer, taskHandler)

	// Reflection
	reflection.Register(grpcServer)

	s.Server = grpcServer
}

func (s *GRPCServer) Start() error {
	if s.isRunning.Load() {
		return fmt.Errorf("server already running")
	}

	s.createServerInstance()

	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.isRunning.Store(true)

	logger.L().Info("gRPC server starting", zap.String("port", s.port))

	go func() {
		if err := s.Server.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
		s.isRunning.Store(false)
	}()

	return nil
}

func (s *GRPCServer) Stop() {
	if !s.isRunning.Load() {
		return
	}
	logger.L().Info("Shutting down gRPC server...")
	s.Server.GracefulStop()
	s.isRunning.Store(false)
}

func (s *GRPCServer) IsRunning() bool {
	return s.isRunning.Load()
}
