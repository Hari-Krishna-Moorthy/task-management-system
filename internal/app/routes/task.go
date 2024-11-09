package routes

import (
	"context"
	"log"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/gofiber/fiber/v2"
)

func setupTaskRoutes(app *fiber.App) error {
	database, err := helpers.GetDatabase(context.TODO())
	if err != nil {
		log.Printf("Error getting database: %v", err)
		panic(err)
	}

	controller := controller.NewTaskController(
		services.GetTaskService(services.GetTaskRepository(database)),
	)

	app.Group("/").
		Get("/tasks", controller.ListTasks).
		Get("/tasks/:id", controller.GetTask).
		Post("/tasks", controller.CreateTask).
		Put("/tasks/:id", controller.UpdateTask).
		Delete("/tasks/:id", controller.DeleteTask)

	return nil
}
