package services

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
)

type TaskService struct {
	taskRepository *TaskRepository
	validator      *validator.Validate
}

type TaskServiceInterface interface {
	CreateTask(ctx context.Context, req *types.CreateTaskRequest) error
	UpdateTask(ctx context.Context, req *types.UpdateTaskRequest) error
	GetTask(ctx context.Context, req *types.GetTaskRequest) error
	DeleteTask(ctx context.Context, req *types.DeleteTaskRequest) error
	ListTasks(ctx context.Context, req *types.ListTasksRequest) error
}

var taskService *TaskService

func NewTaskService(taskRepository *TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: taskRepository,
		validator:      validator.New(),
	}
}

func GetTaskService(taskRepository *TaskRepository) *TaskService {
	if taskService == nil {
		taskService = NewTaskService(taskRepository)
	}
	return taskService
}
