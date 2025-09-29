package grpc

import (
	"context"
	"fmt"
	"log"
	"time"
	"vado/internal/gui/common"
	"vado/internal/pb/taskpb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClientBoxGRPC(win fyne.Window) fyne.CanvasObject {
	btn := common.NewBtn("Запросить кол-во заданий", nil, func() {
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer func(conn *grpc.ClientConn) {
			_ = conn.Close()
		}(conn) // Закрываем сразу после вызова

		client := taskpb.NewTaskServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.GetAllTasks(ctx, &taskpb.Empty{})
		if err != nil {
			fmt.Printf("could not get tasks: %v\n", err)
			return
		}
		for _, t := range resp.Tasks {
			fmt.Printf("Task: %d %s (%s) completed=%v\n",
				t.Id, t.Name, t.Description, t.Completed)
		}

		dialog.ShowInformation("Ответ gRPC", fmt.Sprintf("Заданий: %d", len(resp.Tasks)), win)
	})

	return btn
}
