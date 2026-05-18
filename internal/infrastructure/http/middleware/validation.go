package middleware

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationMiddleware struct {
	validator *validator.Validate
}

func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator.New(),
	}
}

func (v *ValidationMiddleware) ValidateStruct(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *ValidationMiddleware) Validate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors []string
		for _, fieldErr := range validationErrors {
			errors = append(errors, fmt.Sprintf("%s: %s", fieldErr.Field(), fieldErr.Tag()))
		}
		return strings.Join(errors, ", ")
	}
	return err.Error()
}