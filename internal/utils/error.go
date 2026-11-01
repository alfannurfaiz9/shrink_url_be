package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidatorError(err error) string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, v := range ve {
			field := v.Field()

			switch v.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", field)
			case "email":
				return "invalid email format"
			case "min":
				return fmt.Sprintf("%s must be at least %s characters", field, v.Param())
			case "max":
				return fmt.Sprintf("%s must be at most %s characters", field, v.Param())
			default:
				return fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return ""
}
