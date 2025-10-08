package grpc

import (
	"context"
	"fmt"
	"time"
	"vado/internal/constant"
	"vado/internal/gui/common"
	"vado/internal/pb/taskpb"
	"vado/pkg/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientBoxGRPC(win fyne.Window) fyne.CanvasObject {
	btn := common.NewBtn("Запросить кол-во заданий", nil, func() {

		conn, err := grpc.NewClient(fmt.Sprintf("localhost%s", constant.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.L().Error("gRPC connection failed", zap.Error(err))
			dialog.ShowError(err, win)

		}
		defer func(conn *grpc.ClientConn) {
			if closeErr := conn.Close(); closeErr != nil {
				logger.L().Warn("Failed to close connection", zap.Error(closeErr))
			}
		}(conn)

		client := taskpb.NewTaskServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.GetAllTasks(ctx, nil)
		if err != nil {
			dialog.ShowError(err, win)
			logger.L().Error("gRPC request failed", zap.Error(err))
			return
		}

		dialog.ShowInformation("Ответ gRPC", fmt.Sprintf("Заданий: %d", len(resp.Tasks)), win)
	})

	return btn
}
