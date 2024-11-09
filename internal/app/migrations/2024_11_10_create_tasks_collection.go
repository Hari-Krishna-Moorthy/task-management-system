package migrations

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func createTasksCollection(client *mongo.Client) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(config.GetConfig().Database.DBName)

	// Check if the collection already exists
	exists, err := collectionExists(db, (&models.Task{}).GetCollectionsName())
	if err != nil {
		return err
	}

	if exists {
		log.Println("Collection 'task' already exists, skipping creation.")
		return nil
	}

	// Create the collection
	err = db.CreateCollection(ctx, (&models.Task{}).GetCollectionsName())
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
		return err
	}

	log.Println("Tasks collection and indexes created successfully.")
	return nil
}
