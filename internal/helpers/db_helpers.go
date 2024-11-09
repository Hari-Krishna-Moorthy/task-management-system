package helpers

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newMangoDBClient(ctx context.Context) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.GetConfig().Database.URI)
	clientOptions.SetServerAPIOptions(serverAPI)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetMaxPoolSize(config.GetConfig().Database.MaxPoolSize)
	clientOptions.SetMaxConnIdleTime(time.Duration(config.GetConfig().Database.MaxConnIdleTime))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to MongoDB!")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
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
