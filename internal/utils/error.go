package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidatorError(err error) map[string]string {
	var ve validator.ValidationErrors
	errorMessages := make(map[string]string)

	if errors.As(err, &ve) {
		for _, v := range ve {
			field := v.Field()

			switch v.Tag() {
			case "required":
				errorMessages[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorMessages[field] = "invalid email format"
			case "min":
				errorMessages[field] = fmt.Sprintf("%s must be at least %s characters", field, v.Param())
			case "max":
				errorMessages[field] = fmt.Sprintf("%s must be at most %s characters", field, v.Param())
			default:
				errorMessages[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errorMessages
}
