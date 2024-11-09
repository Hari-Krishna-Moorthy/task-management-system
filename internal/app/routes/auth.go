package routes

import (
	"context"
	"log"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

func setupAuthRoutes(app *fiber.App) error {
	database, err := helpers.GetDatabase(context.TODO())
	if err != nil {
		log.Printf("Error getting database: %v", err)
		panic(err)
	}
	controller := controller.NewAuthController(
		services.InitializeAuthService(context.TODO(), database),
	)

	app.Group("/").
		Post("/signup", controller.SignUp).
		Post("/signin", controller.SignIn).
		Post("/signout", controller.SignOut)

	return nil
}
