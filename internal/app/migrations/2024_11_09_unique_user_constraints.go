package migrations

import (
	"context"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// addUniqueUserConstraints applies unique constraints on username, email, and id fields.
func addUniqueUserConstraints(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := client.Database(config.GetConfig().Database.DBName).Collection((&models.User{}).GetCollectionsName())

	indexes := []mongo.IndexModel{
		{Keys: bson.M{"username": 1}, Options: options.Index().SetUnique(true).SetName("unique_username")},
		{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true).SetName("unique_email")},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, indexes)
	return err
}
