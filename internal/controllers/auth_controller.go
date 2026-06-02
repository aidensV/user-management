package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/internal/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Login handles user login
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=models.LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	ipAddress := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	response, err := c.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), "AUTH_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, response, "Login successful")
}

// Logout handles user logout
// @Summary User logout
// @Tags Auth
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := c.authService.Logout(token); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to logout", "LOGOUT_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Logout successful")
}

// GetMe returns current authenticated user info
// @Summary Get current user info
// @Tags Auth
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/me [get]
func (c *AuthController) GetMe(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	user, err := c.authService.GetUserByToken(token)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error(), "INVALID_TOKEN")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, user, "User info retrieved")
}
