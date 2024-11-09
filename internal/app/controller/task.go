package controller

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	"github.com/gofiber/fiber/v2"
)

type TaskController struct {
	TaskService *services.TaskService
}

type TaskControllerInterface interface {
	CreateTask(c *fiber.Ctx) error
	UpdateTask(c *fiber.Ctx) error
	GetTask(c *fiber.Ctx) error
	DeleteTask(c *fiber.Ctx) error
	ListTasks(c *fiber.Ctx) error
}

var taskController TaskControllerInterface

func NewTaskController(taskService *services.TaskService) TaskControllerInterface {
	return &TaskController{
		TaskService: taskService,
	}
}

func GetTaskController(taskService *services.TaskService) TaskControllerInterface {
	if taskController == nil {
		taskController = NewTaskController(taskService)
	}
	return taskController
}

func (t *TaskController) CreateTask(c *fiber.Ctx) error {
	return nil
}

func (t *TaskController) UpdateTask(c *fiber.Ctx) error {
	return nil
}

func (t *TaskController) GetTask(c *fiber.Ctx) error {
	return nil
}

func (t *TaskController) DeleteTask(c *fiber.Ctx) error {
	return nil
}

func (t *TaskController) ListTasks(c *fiber.Ctx) error {
	return nil
}
