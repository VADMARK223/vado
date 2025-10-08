package admin

import (
	"vado/internal/domain/user"
	"vado/internal/gui/common"
	"vado/internal/server/context"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

func NewUserControl(appCtx *context.AppContext, win fyne.Window) fyne.CanvasObject {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Введите имя пользователя")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Введите пароль пользователя")

	createBtn := common.NewBtn("Создать", nil, func() {
		newUser := user.User{Username: usernameEntry.Text, Password: passwordEntry.Text}
		err := appCtx.HTTP.UserService.CreateUser(newUser)
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

	//if util.AutoStartServerGRPC() {
	//	startServerGRPC()
	//}

	return grid
}

//type UserServiceServer struct {
//	userpb.UnimplementedCreateUserServer
//}

/*func startServerGRPC() {
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
}*/
