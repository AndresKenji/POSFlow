package response

import (
	"POSFlowBackend/internal/domain/shared"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// Success sends a successful response with data
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}, message string) {
	Success(c, http.StatusCreated, data, message)
}

// OK sends a 200 OK response
func OK(c *gin.Context, data interface{}, message string) {
	Success(c, http.StatusOK, data, message)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context, message string) {
	c.JSON(http.StatusNoContent, APIResponse{
		Success: true,
		Message: message,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, err error, message string) {
	errorInfo := &ErrorInfo{
		Code:    getErrorCode(err),
		Details: err.Error(),
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errorInfo,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, err error, message string) {
	Error(c, http.StatusBadRequest, err, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, shared.ErrNotFound, message)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, err error, message string) {
	Error(c, http.StatusInternalServerError, err, message)
}

// UnprocessableEntity sends a 422 Unprocessable Entity response
func UnprocessableEntity(c *gin.Context, err error, message string) {
	Error(c, http.StatusUnprocessableEntity, err, message)
}

// HandleError automatically determines the appropriate error response based on the error type
func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		NotFound(c, "Resource not found")
	case errors.Is(err, shared.ErrInvalidInput):
		BadRequest(c, err, "Invalid input provided")
	case errors.Is(err, shared.ErrInsufficientStock):
		UnprocessableEntity(c, err, "Operation failed: insufficient stock")
	case errors.Is(err, shared.ErrInvalidPrice):
		BadRequest(c, err, "Invalid price value")
	case errors.Is(err, shared.ErrInvalidQuantity):
		BadRequest(c, err, "Invalid quantity value")
	case errors.Is(err, shared.ErrOrderNotModifiable):
		UnprocessableEntity(c, err, "Order cannot be modified")
	default:
		InternalServerError(c, err, "Internal server error occurred")
	}
}

// getErrorCode returns a machine-readable error code
func getErrorCode(err error) string {
	switch {
	case errors.Is(err, shared.ErrNotFound):
		return "NOT_FOUND"
	case errors.Is(err, shared.ErrInvalidInput):
		return "INVALID_INPUT"
	case errors.Is(err, shared.ErrInsufficientStock):
		return "INSUFFICIENT_STOCK"
	case errors.Is(err, shared.ErrInvalidPrice):
		return "INVALID_PRICE"
	case errors.Is(err, shared.ErrInvalidQuantity):
		return "INVALID_QUANTITY"
	case errors.Is(err, shared.ErrOrderNotModifiable):
		return "ORDER_NOT_MODIFIABLE"
	default:
		return "INTERNAL_ERROR"
	}
}
