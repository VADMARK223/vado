package grpc

import (
	"context"
	"vado/internal/domain/task"
	"vado/internal/pb/taskpb"
)

// TaskGRPCHandler реализует gRPC сервис для задач
type TaskGRPCHandler struct {
	taskpb.UnimplementedTaskServiceServer
	service task.ITaskService
}

func NewTaskGRPCHandler(service task.ITaskService) *TaskGRPCHandler {
	return &TaskGRPCHandler{service: service}
}

func (h *TaskGRPCHandler) GetAllTasks(ctx context.Context, req *taskpb.Empty) (*taskpb.TaskList, error) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return DomainTaskListToPB(tasks), nil
}

func (h *TaskGRPCHandler) CreateTask(ctx context.Context, req *taskpb.Task) (*taskpb.Empty, error) {
	domainTask := PBTaskToDomain(req)
	err := h.service.CreateTask(*domainTask)
	if err != nil {
		return nil, err
	}
	return &taskpb.Empty{}, nil
}

func (h *TaskGRPCHandler) GetTaskByID(ctx context.Context, req *taskpb.TaskID) (*taskpb.Task, error) {
	domainTask, err := h.service.GetTaskByID(int(req.Id))
	if err != nil {
		return nil, err
	}
	return DomainTaskToPB(domainTask), nil
}

func (h *TaskGRPCHandler) DeleteTask(ctx context.Context, req *taskpb.TaskID) (*taskpb.Empty, error) {
	err := h.service.DeleteTask(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &taskpb.Empty{}, nil
}

func (h *TaskGRPCHandler) DeleteAllTasks(ctx context.Context, req *taskpb.Empty) (*taskpb.Empty, error) {
	h.service.DeleteAllTasks()
	return &taskpb.Empty{}, nil
}
