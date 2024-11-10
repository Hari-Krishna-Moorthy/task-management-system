package services

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	taskCollections *mongo.Collection
}

type TaskRepositoryInterface interface {
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, id string, userID string) error
	GetTask(ctx context.Context, id string, userID string) (*models.Task, error)
	ListTasks(ctx context.Context, userID string) ([]*models.Task, error)
}

var taskRepository TaskRepositoryInterface

func NewTaskRepository(database *mongo.Database) TaskRepositoryInterface {
	return &TaskRepository{
		taskCollections: database.Collection((&models.Task{}).GetCollectionsName()),
	}
}

func GetTaskRepository(database *mongo.Database) TaskRepositoryInterface {
	if taskRepository == nil {
		taskRepository = NewTaskRepository(database)
	}
	return taskRepository
}

func (taskRepository *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	insertOneResult, err := taskRepository.taskCollections.InsertOne(ctx, task)

	if err != nil {
		return err
	}

	log.Println("\nTask created successfully with id: ", insertOneResult.InsertedID)
	return nil
}

func (taskRepository *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	task.UpdatedAt = time.Now()

	update := bson.M{"updated_at": task.UpdatedAt}

	if task.Completed {
		update["completed"] = task.Completed
	}
	if task.Title != utils.EmptyString {
		update["title"] = task.Title
	}
	if task.Description != utils.EmptyString {
		update["description"] = task.Description
	}

	if task.Status != utils.NumbeZero {
		update["status"] = task.Status
	}

	if task.DueDate != utils.DeleteAtZeroTime {
		update["due_date"] = task.DueDate
	}

	if task.DeletedAt != utils.DeleteAtZeroTime {
		update[utils.DeletedAt] = task.DeletedAt
	}

	updateOneResult, err := taskRepository.taskCollections.UpdateByID(ctx, task.ID, bson.M{"$set": update})

	if err != nil {
		return err
	}

	if updateOneResult.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	log.Println("\nTask updated successfully : ", updateOneResult)
	return nil
}

func (taskRepository *TaskRepository) GetTask(ctx context.Context, id string, userID string) (*models.Task, error) {
	var task models.Task
	err := taskRepository.taskCollections.FindOne(ctx, bson.M{utils.IdKey: id, utils.UserIDKey: userID, utils.DeletedAt: utils.DeleteAtZeroTime}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		log.Println("Task not found for id: ", id)
		return nil, err
	}
	return &task, err
}

func (taskRepository *TaskRepository) DeleteTask(ctx context.Context, id string, userID string) error {
	return taskRepository.UpdateTask(ctx, &models.Task{ID: id, UserID: userID, DeletedAt: time.Now().UTC()})
}

func (taskRepository *TaskRepository) ListTasks(ctx context.Context, userID string) ([]*models.Task, error) {
	var tasks []*models.Task
	cursor, err := taskRepository.taskCollections.Find(ctx, bson.M{utils.UserIDKey: userID, utils.DeletedAt: utils.DeleteAtZeroTime})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}
