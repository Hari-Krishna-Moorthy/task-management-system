package migrations

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func collectionExists(db *mongo.Database, collectionName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return false, err
	}

	for _, collName := range collections {
		if collName == collectionName {
			return true, nil
		}
	}
	return false, nil
}

func createUsersCollection(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(config.GetConfig().Database.DBName)

	// Check if the collection already exists
	exists, err := collectionExists(db, (&models.User{}).GetCollectionsName())
	if err != nil {
		return err
	}

	if exists {
		log.Println("Collection 'users' already exists, skipping creation.")
		return nil
	}

	// Create the collection
	err = db.CreateCollection(ctx, (&models.User{}).GetCollectionsName())
	if err != nil {
		log.Printf("Failed to create collection: %v", err)
		return err
	}

	log.Println("Users collection and indexes created successfully.")
	return nil
}
