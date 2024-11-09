package models

import (
	"time"
)

// User represents a user model in the database
type User struct {
	ID       string `json:"id" bson:"_id"` // Unique user ID
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`

	CreatedAt time.Time `json:"created_at" bson:"created_at" `
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" `
	DeletedAt time.Time `json:"deleted_at" bson:"deleted_at" `
}

func (u *User) GetCollectionsName() string {
	return "users"
}
