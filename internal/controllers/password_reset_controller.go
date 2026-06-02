package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"user-management/internal/services"
	"user-management/internal/utils"
)

type PasswordResetController struct {
	resetService *services.PasswordResetService
}

func NewPasswordResetController() *PasswordResetController {
	return &PasswordResetController{
		resetService: services.NewPasswordResetService(),
	}
}

// ForgotPassword handles forgot password request
// @Summary Forgot password
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/forgot-password [post]
func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	token, err := c.resetService.GenerateResetToken(req.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GENERATE_TOKEN_FAILED")
		return
	}

	// TODO: Send email with reset link
	// For now, just return the token (in production, send via email)
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"message": "If email exists, reset token has been sent",
		"token":   token, // Remove this in production, just for testing
	}, "Password reset token generated")
}

// ResetPassword handles password reset
// @Summary Reset password
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Token and new password"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/reset-password [post]
func (c *PasswordResetController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	if len(req.NewPassword) < 6 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Password must be at least 6 characters", "INVALID_PASSWORD")
		return
	}

	if err := c.resetService.ResetPassword(req.Token, req.NewPassword); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "RESET_PASSWORD_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Password reset successfully. Please login with your new password")
}

// ValidateResetToken validates a reset token
// @Summary Validate reset token
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body ValidateTokenRequest true "Token"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/auth/validate-reset-token [post]
func (c *PasswordResetController) ValidateResetToken(ctx *gin.Context) {
	var req ValidateTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	resetToken, err := c.resetService.ValidateResetToken(req.Token)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "INVALID_TOKEN")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"valid": resetToken != nil,
		"user_id": func() string {
			if resetToken != nil {
				return resetToken.UserID.String()
			}
			return ""
		}(),
	}, "Token validation result")
}

// Request DTOs
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
