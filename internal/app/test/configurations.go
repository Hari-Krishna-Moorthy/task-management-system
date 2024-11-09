package testhelpers

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/migrations"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() { //nolint:gochecknoinits
	config.LoadConfig("test") // nolint:errcheck,nolintlint
}

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
		database.Collection(collection).DeleteMany(context.Background(), bson.M{}) //nolint:errcheck,nolintlint
	}
}

func GenerateToken(user *models.User) (string, error) {
	var jwtSecret = []byte(config.GetConfig().Auth.JWTSecret)

	if len(jwtSecret) == 0 {
		jwtSecret = []byte(utils.JWT_DEFAULT_SECRET)
		log.Println("JWT secret not found in environment, using default")
	}

	log.Printf("Generating token for user: %s", user.ID)
	claims := &types.JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: time.Now().Unix(),
		ExpireAt:  time.Now().Add(utils.JWT_TOKEN_EXPIRY).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}
	log.Println("Token generated successfully")
	return signedToken, nil
}
