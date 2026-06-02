package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"user-management/internal/utils"
)

// RequireRole checks if the authenticated user has any of the allowed roles
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (set by auth middleware)
		role, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid role format", "INVALID_ROLE")
			c.Abort()
			return
		}

		// Check if user's role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions. Required role: "+stringsJoin(allowedRoles, ", "), "FORBIDDEN")
		c.Abort()
	}
}

// RequireAdmin is a shortcut for RequireRole("admin")
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireAdminOrSupervisor is a shortcut for admin or supervisor
func RequireAdminOrSupervisor() gin.HandlerFunc {
	return RequireRole("admin", "supervisor")
}

// RequireAdminOrOperator is a shortcut for admin or warehouse_operator
func RequireAdminOrOperator() gin.HandlerFunc {
	return RequireRole("admin", "warehouse_operator")
}

// Helper function
func stringsJoin(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
