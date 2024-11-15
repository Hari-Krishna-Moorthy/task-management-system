// internal/app/controller/auth.go
package controller

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthServiceInterface
}

type AuthControllerInterface interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
}

var authController AuthControllerInterface

func NewAuthController(authService services.AuthServiceInterface) AuthControllerInterface {
	return &AuthController{
		AuthService: authService,
	}
}

func GetAuthController(authService services.AuthServiceInterface) AuthControllerInterface {
	if authController == nil {
		authController = NewAuthController(authService)
	}
	return authController
}

func (authController *AuthController) SignUp(c *fiber.Ctx) error {
	var req types.SignUpRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	return authController.AuthService.SignUp(c, &req)
}

func (authController *AuthController) SignIn(c *fiber.Ctx) error {
	var req types.SignInRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	return authController.AuthService.SignIn(c, &req)
}

func (authController *AuthController) SignOut(c *fiber.Ctx) error {
	return authController.AuthService.SignOut(c)
}
