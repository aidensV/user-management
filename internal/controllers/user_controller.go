package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-management/internal/models"
	"user-management/internal/services"
	"user-management/internal/utils"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// CreateUser handles creating a new user (admin only)
// @Summary Create user
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	// Get current user ID from context (set by auth middleware)
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}

	userID, err := uuid.Parse(currentUserID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	user, err := c.userService.CreateUser(&req, userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "CREATE_USER_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, user, "User created successfully")
}

// GetUser retrieves a user by ID
// @Summary Get user by ID
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/{userId} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	userIDParam := ctx.Param("userId") // ← ganti dari "id" jadi "userId"
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_ID")
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_USER_FAILED")
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", "USER_NOT_FOUND")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, user, "User retrieved")
}

// GetAllUsers retrieves all users with pagination
// @Summary Get all users
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param search query string false "Search by email or name"
// @Success 200 {object} utils.Response{data=[]models.UserResponse}
// @Router /api/v1/users [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	search := ctx.DefaultQuery("search", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := c.userService.GetAllUsers(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_USERS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}, "Users retrieved")
}

// UpdateUser updates an existing user
// @Summary Update user
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Param request body models.UpdateUserRequest true "Update data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/{userId} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userIDParam := ctx.Param("userId") // ← ganti dari "id" jadi "userId"
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_ID")
		return
	}

	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}

	updatedBy, err := uuid.Parse(currentUserID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	user, err := c.userService.UpdateUser(userID, &req, updatedBy)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "UPDATE_USER_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, user, "User updated successfully")
}

// DeleteUser soft deletes a user
// @Summary Delete user (soft delete)
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/{userId} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIDParam := ctx.Param("userId") // ← ganti dari "id" jadi "userId"
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_ID")
		return
	}

	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}

	deletedBy, err := uuid.Parse(currentUserID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	if err := c.userService.DeleteUser(userID, deletedBy); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "DELETE_USER_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "User deleted successfully")
}
