package middleware

import (
	"POSFlowBackend/internal/domain/shared"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Map domain errors to HTTP status codes
			statusCode := mapErrorToStatusCode(err)

			c.JSON(statusCode, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}

func mapErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, shared.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrInsufficientStock):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrInvalidPrice):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrInvalidQuantity):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrOrderNotModifiable):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
