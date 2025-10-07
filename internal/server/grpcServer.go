package server

import (
	"sync"

	"google.golang.org/grpc"
)

var (
	GrpcServer      *grpc.Server
	GrpcServerMutex sync.Mutex
)
