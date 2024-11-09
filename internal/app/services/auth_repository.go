package services

import (
	"context"
	"errors"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errUserExists = errors.New("user already exists")
	errUserNotFound = errors.New("user not found")
)

type AuthRepository struct {
	authCollections *mongo.Collection
}

var authRepository *AuthRepository

func NewAuthRepository(database *mongo.Database) *AuthRepository {
	return &AuthRepository{
		authCollections: database.Collection((&models.User{}).GetCollectionsName()),
	}
}

func GetAuthRepository(database *mongo.Database) *AuthRepository {
	if authRepository == nil {
		authRepository = NewAuthRepository(database)
	}
	return authRepository
}

func (authRepository *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := authRepository.authCollections.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errUserNotFound
	}
	return &user, err
}

func (authRepository *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := authRepository.authCollections.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errUserNotFound
	}
	return &user, err
}

func (authRepository *AuthRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := authRepository.authCollections.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errUserNotFound
	}
	return &user, err
}

func (authRepository *AuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now().UTC()
	_, err := authRepository.authCollections.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (authRepository *AuthRepository) DeleteUser(ctx context.Context, id string) error {
	return authRepository.UpdateUser(ctx, &models.User{ID: id, DeletedAt: time.Now().UTC()})
}

func (authRepository *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := authRepository.checkUserExists(ctx, user); err != nil {
		return err
	}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	_, err := authRepository.authCollections.InsertOne(ctx, user)
	return err
}

func (authRepository *AuthRepository) checkUserExists(ctx context.Context, user *models.User) error {

	filter := bson.M{
		utils.MongoDBFilterOr: []bson.M{
			{utils.Username: user.Username},
			{utils.Email: user.Email},
		},
	}

	var existingUser *models.User
	err := authRepository.authCollections.FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return errUserExists
	} else if err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}
