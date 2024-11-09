package types

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/golang-jwt/jwt"
)

// Auth request types
type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30" validateMsg:"Username is required and must be between 3 and 30 characters"`
	Email    string `json:"email" validate:"required,email" validateMsg:"Valid email is required"`
	Password string `json:"password" validate:"required,min=8,max=128" validateMsg:"Password is required and must be between 8 and 128 characters"`
}

type SignUpResponse struct {
	User    *models.User `json:"user"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required_without=Username" validateMsg:"Email or Username is required, and must be a valid email"`
	Username string `json:"username" validate:"required_without=Email" validateMsg:"Username or Email is required and must be between 3 and 30 characters"`
	Password string `json:"password" validate:"required,min=8,max=128" validateMsg:"Password is required and must be between 8 and 128 characters"`
}

type SignInResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SignOutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// JWT token
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
	ExpireAt  int64  `json:"expire_at"`
	jwt.StandardClaims
}
