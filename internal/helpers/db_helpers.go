package helpers

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMangoDBClient(ctx context.Context) (*mongo.Client, error) {

	commandMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, event *event.CommandStartedEvent) {
			log.Printf(string(utils.BlueColor)+"\nMongoDB `Query Started`:\n%+v\n%v", event.Command, string(utils.ResetColor))
		},
		Succeeded: func(_ context.Context, event *event.CommandSucceededEvent) {
			log.Printf(string(utils.GreenColor)+"\nMongoDB Query Succeeded:\nCommand Name: %v\nDuration: %v\n%v",
				event.CommandName, event.Duration, string(utils.ResetColor))
		},
		Failed: func(_ context.Context, event *event.CommandFailedEvent) {
			log.Printf(string(utils.RedColor)+"\nMongoDB Query Failed:\nCommand Name: %v\nDuration: %v\nFailure: %v\n%v",
				event.CommandName, event.Duration, event.Failure, string(utils.ResetColor))
		},
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.GetConfig().Database.URI)
	clientOptions.SetServerAPIOptions(serverAPI)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetMaxPoolSize(config.GetConfig().Database.MaxPoolSize)
	clientOptions.SetMaxConnIdleTime(time.Duration(config.GetConfig().Database.MaxConnIdleTime))
	clientOptions.SetMonitor(commandMonitor)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %w", err)
		return nil, err
	}
	return client, nil
}

func GetClient(ctx context.Context) (*mongo.Client, error) {
	return newMangoDBClient(ctx)
}

func GetDatabase(ctx context.Context) (*mongo.Database, error) {
	client, err := newMangoDBClient(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(config.GetConfig().Database.DBName), nil
}

func GetCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	db, err := GetDatabase(ctx)
	if err != nil {
		return nil, err
	}

	return db.Collection(collectionName), nil
}
