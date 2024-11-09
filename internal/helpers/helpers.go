package helpers

import (
	"errors"
	"log"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func FormateValidationError(err error) []*types.RequestError {
	if err == nil {
		return nil
	}
	var errors []*types.RequestError
	for _, err := range err.(validator.ValidationErrors) {
		var el types.RequestError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors = append(errors, &el)
	}

	return errors
}

func GetUserDataFromToken(tokenString string) (string, error) {
	jwtSecret := []byte(config.GetConfig().Auth.JWTSecret)
	if len(jwtSecret) == 0 {
		jwtSecret = []byte(utils.JWT_DEFAULT_SECRET)
		log.Println("JWT secret not found in environment, using default")
	}

	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return "", err
	}

	claims, ok := token.Claims.(*types.JWTClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.UserID, nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
