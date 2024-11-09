package test_helpers

import (
	"context"
	"log"
	"testing"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/migrations"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	config.LoadConfig("test") // nolint:errcheck
}

var ()

func initializesDB(t *testing.T) *mongo.Client {

	client, err := helpers.GetClient(context.Background())
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return client
}

func GetTestClient() *mongo.Client {
	return initializesDB(nil)
}

func GetTestDatabase() *mongo.Database {
	client := GetTestClient()
	return client.Database(config.GetConfig().Database.DBName)
}

func initMigrations(t *testing.T) {
	if err := migrations.ApplyMigrations(GetTestClient()); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
}

func initTeardown(t *testing.T) {
	if GetTestClient() != nil {
		if err := GetTestClient().Disconnect(context.Background()); err != nil {
			t.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}
}

func TestSetup(t *testing.T) {
	env := "test"
	err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	initializesDB(t)
	initMigrations(t)
	defer initTeardown(t)
}

func TruncatesCollection() {
	database := GetTestDatabase()
	collectionNames := []string{"tasks", "users"}
	for _, collection := range collectionNames {
		database.Collection(collection).DeleteMany(context.Background(), bson.M{})
	}
}
