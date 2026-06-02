package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"user-management/internal/services"
	"user-management/internal/utils"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		authService: services.NewAuthService(),
	}
}

// Authenticate verifies the JWT token and sets user info in context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header", "MISSING_TOKEN")
			c.Abort()
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization format. Use Bearer <token>", "INVALID_TOKEN_FORMAT")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token: "+err.Error(), "INVALID_TOKEN")
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("name", claims.Name)

		c.Next()
	}
}

// OptionalAuth attempts to authenticate but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token := parts[1]
				claims, err := m.authService.ValidateToken(token)
				if err == nil {
					c.Set("user_id", claims.UserID.String())
					c.Set("email", claims.Email)
					c.Set("role", claims.Role)
					c.Set("name", claims.Name)
					c.Set("authenticated", true)
				}
			}
		}
		c.Next()
	}
}
