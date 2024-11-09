package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/migrations"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/routes"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/gofiber/fiber/v2"
)

func init() {

	// Enable logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}


	// Load the appropriate configuration
	err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

}

func main() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Printf("Panic occurred: %v", err)
	// 		log.Fatal(err)
	// 	}
	// }()

	// go func() {

	// Initialize the database
	ctx := context.Background()
	databaseClient, err := helpers.GetClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Migrate the database
	if err := migrations.ApplyMigrations(databaseClient); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	app := fiber.New()

	if err := routes.SetupRoutes(app); err != nil {
		log.Fatal(err)
		panic(err)
	}
	port := config.GetConfig().Server.Port

	log.Printf("Starting server on port %d...", port)
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
	// }()
}
