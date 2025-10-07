package admin

import (
	"fmt"
	"log"
	"net"
	constant2 "vado/internal/constant"
	"vado/internal/gui/common"
	"vado/internal/model"
	"vado/internal/pb/userpb"
	"vado/internal/server"
	"vado/internal/util"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewUserControl(appCtx *util.AppContext, win fyne.Window) fyne.CanvasObject {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Введите пароль пользователя")

	createBtn := common.NewBtn("Создать", nil, func() {
		newUser := model.User{Username: usernameEntry.Text, Password: passwordEntry.Text}
		err := appCtx.HttpContext.UserService.CreateUser(newUser)
		if err != nil {
			dialog.ShowError(err, win)
			logger.L().Error("create user failed", zap.Error(err))
			return
		}
	})

	grid := container.NewGridWithColumns(3,
		usernameEntry,
		passwordEntry,
		createBtn,
	)

	fmt.Println(server.GrpcServer)
	if util.AutoStartServerGRPC() {
		startServerGRPC()
	}

	return grid
}

type UserServiceServer struct {
	userpb.UnimplementedCreateUserServer
}

func startServerGRPC() {
	fmt.Println("GRPC server started.")
	server.GrpcServerMutex.Lock()
	defer server.GrpcServerMutex.Unlock()

	if server.GrpcServer != nil {
		fmt.Println("Server not running")
		return
	}

	var err error
	listener, err := net.Listen("tcp", constant2.GrpcPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	server.GrpcServer = grpc.NewServer()
	srv := &UserServiceServer{}
	userpb.RegisterCreateUserServer(server.GrpcServer, srv)
	reflection.Register(server.GrpcServer)

	go func() {
		logger.L().Info(fmt.Sprintf("gRPC-server started on %s", constant2.GrpcPort))

		if err := server.GrpcServer.Serve(listener); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()
}
