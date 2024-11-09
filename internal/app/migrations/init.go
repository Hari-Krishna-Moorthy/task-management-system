package migrations

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MigrationID string             `bson:"migration_id"` // Unique ID for each migration
	Description string             `bson:"description"`
	AppliedAt   time.Time          `bson:"applied_at"`
}

// MigrationList holds all migrations in order of application
var MigrationList = []struct {
	ID          string
	Description string
	Apply       func(*mongo.Client) error
}{
	{
		ID:          "2024_11_10_create_users_collection",
		Description: "Create users collection with _id, username, email, password, created_at, updated_at fields",
		Apply:       createUsersCollection,
	},
	{
		ID:          "2024_11_09_unique_user_constraints",
		Description: "Add unique constraints on username, email, and id fields in users collection",
		Apply:       addUniqueUserConstraints,
	},
}
