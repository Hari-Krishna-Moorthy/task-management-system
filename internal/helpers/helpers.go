package helpers

import (
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/go-playground/validator/v10"
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
