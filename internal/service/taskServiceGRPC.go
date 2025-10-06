package service

import (
	"context"
	"vado/internal/model"
	"vado/internal/pb/taskpb"
)

type TaskServiceGRPC struct {
	taskpb.UnimplementedTaskServiceServer
	Service ITaskService
}

func NewTaskServiceGRPC(s ITaskService) *TaskServiceGRPC {
	return &TaskServiceGRPC{Service: s}
}

func (s *TaskServiceGRPC) GetAllTasks(ctx context.Context, _ *taskpb.Empty) (*taskpb.TaskList, error) {
	tl, err := s.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return toProtoTaskList(tl), nil
}

func (s *TaskServiceGRPC) CreateTask(ctx context.Context, t *taskpb.Task) (*taskpb.Empty, error) {
	err := s.Service.CreateTask(fromProtoTask(t))
	if err != nil {
		return nil, err
	}
	return &taskpb.Empty{}, nil
}

func (s *TaskServiceGRPC) GetTaskByID(ctx context.Context, id *taskpb.TaskID) (*taskpb.Task, error) {
	task, err := s.Service.GetTaskByID(int(id.Id))
	if err != nil {
		return nil, err
	}
	return toProtoTask(task), nil
}

func (s *TaskServiceGRPC) DeleteTask(ctx context.Context, id *taskpb.TaskID) (*taskpb.Empty, error) {
	err := s.Service.DeleteTask(int(id.Id))
	if err != nil {
		return nil, err
	}
	return &taskpb.Empty{}, nil
}

func (s *TaskServiceGRPC) DeleteAllTasks(ctx context.Context, _ *taskpb.Empty) (*taskpb.Empty, error) {
	s.Service.DeleteAllTasks()
	return &taskpb.Empty{}, nil
}

// toProtoTask - конвертация model.Task в taskpb.Task
func toProtoTask(t *model.Task) *taskpb.Task {
	return &taskpb.Task{
		Id:          int32(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

// fromProtoTask — конвертация taskpb.Task в model.Task
func fromProtoTask(pt *taskpb.Task) model.Task {
	return model.Task{
		ID:          int(pt.Id),
		Name:        pt.Name,
		Description: pt.Description,
		Completed:   pt.Completed,
	}
}

// toProtoTaskList — конвертация model.TaskList в taskpb.TaskList
func toProtoTaskList(tl model.TaskList) *taskpb.TaskList {
	tasks := make([]*taskpb.Task, len(tl.Tasks))
	for i, t := range tl.Tasks {
		tasks[i] = toProtoTask(&t)
	}
	return &taskpb.TaskList{Tasks: tasks}
}

// fromProtoTaskList — конвертация taskpb.TaskList в model.TaskList
func fromProtoTaskList(ptl *taskpb.TaskList) model.TaskList {
	tasks := make([]model.Task, len(ptl.Tasks))
	for i, t := range ptl.Tasks {
		tasks[i] = fromProtoTask(t)
	}
	return model.TaskList{Tasks: tasks}
}

func toModelTask(t *taskpb.Task) model.Task {
	return model.Task{
		ID:          int(t.Id),
		Name:        t.Name,
		Description: t.Description,
		Completed:   t.Completed,
	}
}
