package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"user-management/internal/utils"
)

// Recovery middleware handles panics and returns a JSON error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
				c.Abort()
			}
		}()
		c.Next()
	}
}
