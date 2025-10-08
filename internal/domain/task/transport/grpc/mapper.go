package grpc

import (
	"vado/internal/domain/task"
	"vado/internal/pb/taskpb"
)

// DomainTaskToPB конвертирует domain модель в protobuf
func DomainTaskToPB(t *task.Task) *taskpb.Task {
	if t == nil {
		return nil
	}
	return &taskpb.Task{
		Id:          int32(t.ID),
		Name:        t.Name,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

// PBTaskToDomain конвертирует protobuf в domain модель
func PBTaskToDomain(pb *taskpb.Task) *task.Task {
	if pb == nil {
		return nil
	}
	return &task.Task{
		ID:          int(pb.Id),
		Name:        pb.Name,
		Description: pb.Description,
		Completed:   pb.Completed,
	}
}

// DomainTaskListToPB конвертирует список задач domain в protobuf
func DomainTaskListToPB(list task.TaskList) *taskpb.TaskList {
	pbTasks := make([]*taskpb.Task, 0, len(list.Tasks))
	for _, t := range list.Tasks {
		pbTasks = append(pbTasks, DomainTaskToPB(&t))
	}
	return &taskpb.TaskList{Tasks: pbTasks}
}
