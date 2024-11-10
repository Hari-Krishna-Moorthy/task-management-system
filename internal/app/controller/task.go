package controller

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	TaskService services.TaskServiceInterface
}

type TaskControllerInterface interface {
	CreateTask(c *fiber.Ctx) error
	UpdateTask(c *fiber.Ctx) error
	GetTask(c *fiber.Ctx) error
	DeleteTask(c *fiber.Ctx) error
	ListTasks(c *fiber.Ctx) error
}

var taskController TaskControllerInterface

func NewTaskController(taskService services.TaskServiceInterface) TaskControllerInterface {
	return &TaskController{
		TaskService: taskService,
	}
}

func GetTaskController(taskService services.TaskServiceInterface) TaskControllerInterface {
	if taskController == nil {
		taskController = NewTaskController(taskService)
	}
	return taskController
}

func (t *TaskController) CreateTask(c *fiber.Ctx) error {
	var req types.CreateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	return t.TaskService.CreateTask(c, &req)
}

func (t *TaskController) UpdateTask(c *fiber.Ctx) error {
	var req types.UpdateTaskRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	return t.TaskService.UpdateTask(c, &req)
}

func (t *TaskController) GetTask(c *fiber.Ctx) error {

	return t.TaskService.GetTask(c)
}

func (t *TaskController) DeleteTask(c *fiber.Ctx) error {

	return t.TaskService.DeleteTask(c)
}

func (t *TaskController) ListTasks(c *fiber.Ctx) error {
	var req types.ListTasksRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	return t.TaskService.ListTasks(c, &req)
}
