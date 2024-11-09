package services

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	taskCollections *mongo.Collection
}

var taskRepository *TaskRepository

func NewTaskRepository(database *mongo.Database) *TaskRepository {
	return &TaskRepository{
		taskCollections: database.Collection((&models.Task{}).GetCollectionsName()),
	}
}

func GetTaskRepository(database *mongo.Database) *TaskRepository {
	if taskRepository == nil {
		taskRepository = NewTaskRepository(database)
	}
	return taskRepository
}
