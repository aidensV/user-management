package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, errorMsg string, code string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorMsg,
		Code:    code,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "validation failed",
		"code":    "VALIDATION_ERROR",
		"details": errors,
	})
}
