package services

import (
	"context"
	"log"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	authRepo  *AuthRepository
	validator *validator.Validate
}

var jwtSecret = []byte(config.GetConfig().Auth.JWTSecret)

type AuthServiceInterface interface {
	SignUp(c *fiber.Ctx, req *types.SignUpRequest) error
	SignIn(c *fiber.Ctx, req *types.SignInRequest) error
	SignOut(c *fiber.Ctx) error
}

var authService *AuthService

func InitializeAuthService(ctx context.Context, database *mongo.Database) *AuthService {
	log.Println("Initializing AuthService")
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(utils.JWT_DEFAULT_SECRET)
		log.Println("JWT secret not found in environment, using default")
	}

	return &AuthService{
		authRepo:  GetAuthRepository(database),
		validator: validator.New(),
	}
}

func GetAuthService(ctx context.Context, database *mongo.Database) *AuthService {
	log.Println("Retrieving AuthService")
	if authService == nil {
		log.Println("AuthService not initialized, calling InitializeAuthService")
		authService = InitializeAuthService(ctx, database)
	}
	return authService
}

// GenerateToken creates a JWT token for authenticated users
func (authService *AuthService) generateToken(user *models.User) (string, error) {
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

// SignUp creates a new user in the database
func (authService *AuthService) SignUp(c *fiber.Ctx, req *types.SignUpRequest) error {
	log.Println("SignUp request received")
	response := &types.SignUpResponse{
		Success: false,
	}

	err := authService.validator.Struct(req)
	if err != nil {
		validateionErrors := helpers.FormateValidationError(err)

		if len(validateionErrors) > utils.NumbeZero {
			return c.Status(fiber.StatusNotFound).JSON(validateionErrors)
		}
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := authService.authRepo.CreateUser(c.Context(), user); err != nil {
		log.Printf("Error creating user: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response) // nolint:errcheck

	}
	log.Printf("User %s created successfully", user.ID)

	response.Success = true
	response.Message = "User registered successfully"

	return c.Status(fiber.StatusCreated).JSON(response) // nolint:errcheck
}

func (authService *AuthService) SignIn(c *fiber.Ctx, req *types.SignInRequest) error {
	log.Println("SignIn request received")
	response := &types.SignInResponse{
		Success: false,
	}
	err := authService.validator.Struct(req)
	if err != nil {
		validateionErrors := helpers.FormateValidationError(err)

		if len(validateionErrors) > utils.NumbeZero {
			return c.Status(fiber.StatusNotFound).JSON(validateionErrors)
		}
	}

	var user *models.User

	if req.Email != "" {
		user, err = authService.authRepo.FindUserByEmail(c.Context(), req.Email)
	} else {
		user, err = authService.authRepo.FindUserByUsername(c.Context(), req.Username)
	}

	if err != nil {
		log.Printf("Error finding user: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if !verifyPassword(req.Password, user.Password) {
		log.Println("Invalid email or password")
		response.Message = "Invalid email or password"
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	token, err := authService.generateToken(user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		response.Message = err.Error()
		return c.Status(fiber.StatusNotFound).JSON(response)

	}
	log.Printf("Token generated for user: %s", user.ID)
	cookie := &fiber.Cookie{
		Name:     utils.CookieKeyToken,
		Value:    token,
		Expires:  time.Now().Add(utils.JWT_TOKEN_EXPIRY),
		HTTPOnly: true,
		Domain:   config.GetConfig().Server.AppDomain,
		Path:     "/",
		SameSite: "lax",
		Secure:   true,
		MaxAge:   int(utils.JWT_TOKEN_EXPIRY.Seconds()),
	}
	c.Cookie(cookie)

	response.Token = token
	response.Success = true
	response.Message = "user signed in successfully"
	c.Status(fiber.StatusOK).JSON(response) // nolint:errcheck
	return nil
}

func (authService *AuthService) SignOut(c *fiber.Ctx) error {
	log.Println("SignOut request received")

	if c.Cookies(utils.CookieKeyToken) == "" {
		log.Println("No token found in cookies")
		response := &types.SignOutResponse{
			Success: false,
			Message: "user not signed in",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	c.ClearCookie(utils.CookieKeyToken)
	response := &types.SignOutResponse{
		Success: true,
		Message: "Signed out successfully",
	}
	log.Println("User signed out successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}
