package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func TestSetup(t *testing.T) {
	config.LoadConfig("test") // nolint:errcheck
	uri := config.GetConfig().Database.URI
	fmt.Printf("TEST_DB_URI: %s\n", uri)
	if uri == "" {
		t.Fatal("TEST_DB_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(uri).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())). //nolint:staticcheck
		SetReadConcern(readconcern.Majority()).
		SetMaxPoolSize(config.GetConfig().Database.MaxPoolSize).
		SetMaxConnIdleTime(time.Duration(config.GetConfig().Database.MaxConnIdleTime)).
		SetRetryWrites(config.GetConfig().Database.RetryWrites)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer client.Disconnect(ctx) //nolint:errcheck
}
