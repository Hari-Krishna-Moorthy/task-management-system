package models

import "time"

type Task struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`

	Completed bool      `json:"completed" bson:"completed"`
	Status    string    `json:"status" bson:"status"`
	DueDate   time.Time `json:"due_date" bson:"due_date"`

	UserID string `json:"user_id" bson:"user_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at" `
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" `
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at" `
}
