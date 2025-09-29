package component

import (
	"fmt"
	"log"
	"net"
	"vado/internal/pb/taskpb"
	"vado/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewControlBoxGRPC(s service.ITaskService) fyne.CanvasObject {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		_ = fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	serviceGRPC := service.NewTaskServiceGRPC(s)
	taskpb.RegisterTaskServiceServer(grpcServer, serviceGRPC)

	// включаем рефлексию
	reflection.Register(grpcServer)

	fmt.Println("Listening on port 50051", lis)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return widget.NewLabel("Сервер gGRPC:")
}
