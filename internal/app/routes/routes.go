package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) error {
	// Define your routes here
	if err := setupAuthRoutes(app); err != nil {
		return err
	}

	if err := setupTaskRoutes(app); err != nil {
		return err
	}

	return nil
}
