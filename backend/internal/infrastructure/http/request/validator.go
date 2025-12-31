package request

import (
	"POSFlowBackend/internal/domain/shared"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate binds the request body to the given struct and validates it
func BindAndValidate(c *gin.Context, req interface{}) error {
	// Bind JSON request body
	if err := c.ShouldBindJSON(req); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// BindURI binds URI parameters to the given struct
func BindURI(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindUri(req); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// BindQuery binds query parameters to the given struct
func BindQuery(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindQuery(req); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// formatValidationError formats validation errors into a readable message
func formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Get the first validation error
		if len(validationErrors) > 0 {
			fieldError := validationErrors[0]
			return fmt.Errorf("%w: %s", shared.ErrInvalidInput, formatFieldError(fieldError))
		}
	}

	// Return generic validation error
	return fmt.Errorf("%w: %v", shared.ErrInvalidInput, err)
}

// formatFieldError formats a field validation error
func formatFieldError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", field)
	case "min":
		return fmt.Sprintf("field '%s' must have at least %s items", field, fe.Param())
	case "max":
		return fmt.Sprintf("field '%s' must have at most %s items", field, fe.Param())
	case "gt":
		return fmt.Sprintf("field '%s' must be greater than %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("field '%s' must be greater than or equal to %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("field '%s' must be less than %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("field '%s' must be less than or equal to %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of: %s", field, fe.Param())
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email address", field)
	case "uuid":
		return fmt.Sprintf("field '%s' must be a valid UUID", field)
	default:
		return fmt.Sprintf("field '%s' failed validation: %s", field, fe.Tag())
	}
}

// GetPathParam retrieves a path parameter
func GetPathParam(c *gin.Context, param string) string {
	return c.Param(param)
}

// GetQueryParam retrieves a query parameter with a default value
func GetQueryParam(c *gin.Context, param string, defaultValue string) string {
	if value := c.Query(param); value != "" {
		return value
	}
	return defaultValue
}
