package migrations

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplyMigrations checks and applies any pending migrations.
func ApplyMigrations(client *mongo.Client) error {
	ctx := context.Background()
	migrationsCollection := client.Database(config.GetConfig().Database.DBName).Collection("migrations")

	for _, m := range MigrationList {
		filter := bson.M{"migration_id": m.ID}
		count, err := migrationsCollection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		// Skip this migration if it has already been applied
		if count > 0 {
			log.Printf("Migration %s already applied, skipping.", m.ID)
			continue
		}

		// Apply the migration
		log.Printf("Applying migration: %s - %s", m.ID, m.Description)
		err = m.Apply(client)
		if err != nil {
			log.Printf("Failed to apply migration %s: %v", m.ID, err)
			return err
		}

		// Record the applied migration in the migrations collection
		_, err = migrationsCollection.InsertOne(ctx, Migration{
			MigrationID: m.ID,
			Description: m.Description,
			AppliedAt:   time.Now(),
		})
		if err != nil {
			log.Printf("Failed to record migration %s: %v", m.ID, err)
			return err
		}

		log.Printf("Migration %s applied successfully.", m.ID)
	}

	log.Println("All migrations applied.")
	return nil
}
