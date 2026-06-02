package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-management/internal/repositories"
	"user-management/internal/utils"
)

type LoginHistoryController struct {
	loginHistoryRepo *repositories.LoginHistoryRepository
}

func NewLoginHistoryController() *LoginHistoryController {
	return &LoginHistoryController{
		loginHistoryRepo: repositories.NewLoginHistoryRepository(),
	}
}

// GetMyLoginHistory handles getting current user's login history
// @Summary Get my login history
// @Tags Login History
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/history [get]
func (c *LoginHistoryController) GetMyLoginHistory(ctx *gin.Context) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	histories, err := c.loginHistoryRepo.GetByUserID(userID, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_HISTORY_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, histories, "Login history retrieved")
}

// GetAllLoginHistory handles getting all login history (admin only)
// @Summary Get all login history
// @Tags Login History
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by email"
// @Success 200 {object} utils.Response
// @Router /api/v1/login-history [get]
func (c *LoginHistoryController) GetAllLoginHistory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	search := ctx.DefaultQuery("search", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	histories, total, err := c.loginHistoryRepo.GetAll(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_HISTORY_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"histories":   histories,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}, "Login history retrieved")
}
