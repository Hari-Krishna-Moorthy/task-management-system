package services

import (
	"errors"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/enums"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
)

type TaskService struct {
	taskRepository *TaskRepository
	validator      *validator.Validate
}

type TaskServiceInterface interface {
	CreateTask(c *fiber.Ctx, req *types.CreateTaskRequest) error
	UpdateTask(c *fiber.Ctx, req *types.UpdateTaskRequest) error
	GetTask(c *fiber.Ctx) error
	DeleteTask(c *fiber.Ctx) error
	ListTasks(c *fiber.Ctx, req *types.ListTasksRequest) error
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

func (taskService *TaskService) CreateTask(c *fiber.Ctx, req *types.CreateTaskRequest) error {
	log.Println("Create Task request received", req)
	response := &types.CreateTaskResponse{
		Success: false,
	}

	err := taskService.validator.Struct(req)
	if err != nil {
		validateionErrors := helpers.FormateValidationError(err)

		if len(validateionErrors) > utils.NumbeZero {
			return c.Status(fiber.StatusNotFound).JSON(validateionErrors)
		}
	}

	user, err := validateToken(c)

	if err != nil {
		log.Printf("Error getting user data from token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	dueDate, err := time.Parse(utils.TimeLayout, req.DueDate)

	if err != nil {
		log.Printf("Error parsing due date: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	task := &models.Task{
		ID:          helpers.GenerateUUID(),
		Title:       req.Title,
		Description: req.Description,
		UserID:      user.ID,
		Status:      enums.ToDo,
		DueDate:     dueDate,
	}

	if err := taskService.taskRepository.CreateTask(c.Context(), task); err != nil {
		log.Printf("Error creating task: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response.Success = true
	response.Message = "Task created successfully"
	response.Task = task
	return c.Status(fiber.StatusOK).JSON(response)
}

func (taskService *TaskService) UpdateTask(c *fiber.Ctx, req *types.UpdateTaskRequest) error {
	log.Println("Update Task request received", req)
	response := &types.GetTaskResponse{
		Success: false,
	}

	err := taskService.validator.Struct(req)
	if err != nil {
		validateionErrors := helpers.FormateValidationError(err)

		if len(validateionErrors) > utils.NumbeZero {
			return c.Status(fiber.StatusNotFound).JSON(validateionErrors)
		}
	}
	user, err := validateToken(c)

	if err != nil {
		log.Printf("Error getting user data from token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	id := c.Params("id")

	if req.Status != utils.EmptyString {
		task, err := taskService.taskRepository.GetTask(c.Context(), id, user.ID)

		if err != nil {
			log.Printf("Error getting task: %v", err)
			response.Message = err.Error()
			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		if task.CanMoveStatus(task.Status, enums.TaskStatusFromString(req.Status)) {
			task.ApplyStatusChange(enums.TaskStatusFromString(req.Status))
		}
	}

	task := &models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		UserID:      user.ID,
		Status:      enums.TaskStatusFromString(req.Status),
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse(utils.TimeLayout, req.DueDate)
		if err != nil {
			log.Printf("Error parsing due date: %v", err)
			response.Message = err.Error()
			return c.Status(fiber.StatusNotFound).JSON(response)
		}
		task.DueDate = dueDate
	}

	if err := taskService.taskRepository.UpdateTask(c.Context(), task); err != nil {
		log.Printf("Error updating task: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response.Success = true
	response.Message = "Task updated successfully"
	response.Task = task
	return c.Status(fiber.StatusOK).JSON(response)
}

func (taskService *TaskService) GetTask(c *fiber.Ctx) error {
	log.Println("Create Task request received")
	response := &types.CreateTaskResponse{
		Success: false,
	}
	user, err := validateToken(c)
	if err != nil {
		log.Printf("Error getting user data from token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	id := c.Params("id")

	task, err := taskService.taskRepository.GetTask(c.Context(), id, user.ID)

	if err != nil {
		log.Printf("Error getting task: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response.Success = true
	response.Message = "Task fetched successfully"
	response.Task = task
	return c.Status(fiber.StatusOK).JSON(response)
}

func (taskService *TaskService) DeleteTask(c *fiber.Ctx) error {
	log.Println("Delete Task request received")
	response := &types.DeleteTaskResponse{
		Success: false,
	}

	user, err := validateToken(c)
	if err != nil {
		log.Printf("Error getting user data from token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	id := c.Params("id")

	if err := taskService.taskRepository.DeleteTask(c.Context(), id, user.ID); err != nil {
		log.Printf("Error deleting task: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response.Success = true
	response.Message = "Task deleted successfully"
	return c.Status(fiber.StatusOK).JSON(response)
}

func (taskService *TaskService) ListTasks(c *fiber.Ctx, req *types.ListTasksRequest) error {
	log.Println("List Tasks request received", req)

	response := &types.ListTasksResponse{
		Success: false,
	}

	err := taskService.validator.Struct(req)
	if err != nil {
		validateionErrors := helpers.FormateValidationError(err)

		if len(validateionErrors) > utils.NumbeZero {
			return c.Status(fiber.StatusNotFound).JSON(validateionErrors)
		}
	}

	user, err := validateToken(c)

	if err != nil {
		log.Printf("Error getting user data from token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	tasks, err := taskService.taskRepository.ListTasks(c.Context(), user.ID)

	if err != nil {
		log.Printf("Error getting tasks: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response.Success = true
	response.Message = "Tasks fetched successfully"
	response.Tasks = tasks
	return c.Status(fiber.StatusOK).JSON(response)
}

func validateToken(c *fiber.Ctx) (*models.User, error) {

	userID, err := helpers.GetUserDataFromToken(c.Cookies(utils.CookieKeyToken))
	if err != nil {
		return nil, err
	}

	if userID == utils.EmptyString {
		return nil, errors.New("invalid token")
	}

	user, err := GetAuthRepository(nil).FindUserByID(c.Context(), userID)

	if err != nil || user.ID == utils.EmptyString {
		return nil, err
	}

	return user, nil
}
