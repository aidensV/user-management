package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-management/internal/services"
	"user-management/internal/utils"
)

type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{
		roleService: services.NewRoleService(),
	}
}

// GetAllRoles handles getting all roles
// @Summary Get all roles
// @Tags Roles
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /api/v1/roles [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_ROLES_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, roles, "Roles retrieved")
}

// GetRoleByID handles getting a role by ID
// @Summary Get role by ID
// @Tags Roles
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/roles/{roleId} [get]
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ID")
		return
	}

	role, err := c.roleService.GetRoleByID(roleID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), "ROLE_NOT_FOUND")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, role, "Role retrieved")
}

// CreateRole handles creating a new role
// @Summary Create role
// @Tags Roles
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param request body CreateRoleRequest true "Role data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	role, err := c.roleService.CreateRole(req.Name, req.DisplayName, req.Description)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "CREATE_ROLE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, role, "Role created successfully")
}

// UpdateRole handles updating a role
// @Summary Update role
// @Tags Roles
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Param request body UpdateRoleRequest true "Role data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/roles/{roleId} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ID")
		return
	}

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	role, err := c.roleService.UpdateRole(roleID, req.DisplayName, req.Description, req.IsActive)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "UPDATE_ROLE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, role, "Role updated successfully")
}

// DeleteRole handles deleting a role
// @Summary Delete role
// @Tags Roles
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/roles/{roleId} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ID")
		return
	}

	if err := c.roleService.DeleteRole(roleID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "DELETE_ROLE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Role deleted successfully")
}

// AssignRoleToUser handles assigning a role to a user
// @Summary Assign role to user
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Param request body AssignRoleRequest true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/users/{userId}/roles [post]
func (c *RoleController) AssignRoleToUser(ctx *gin.Context) {
	userIDParam := ctx.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	var req AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	// Get current user ID from context
	currentUserIDStr, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "UNAUTHORIZED")
		return
	}

	assignedBy, err := uuid.Parse(currentUserIDStr.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	if err := c.roleService.AssignRoleToUser(userID, roleID, assignedBy); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "ASSIGN_ROLE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Role assigned successfully")
}

// RemoveRoleFromUser handles removing a role from a user
// @Summary Remove role from user
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/users/{userId}/roles/{roleId} [delete]
func (c *RoleController) RemoveRoleFromUser(ctx *gin.Context) {
	userIDParam := ctx.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	if err := c.roleService.RemoveRoleFromUser(userID, roleID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "REMOVE_ROLE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Role removed successfully")
}

// GetUserRoles handles getting user's roles
// @Summary Get user's roles
// @Tags Users
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param userId path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/users/{userId}/roles [get]
func (c *RoleController) GetUserRoles(ctx *gin.Context) {
	userIDParam := ctx.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	roles, err := c.roleService.GetUserRoles(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_ROLES_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, roles, "User roles retrieved")
}

// Request DTOs
type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=50"`
	DisplayName string  `json:"display_name" binding:"required,min=3,max=100"`
	Description *string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	DisplayName string  `json:"display_name" binding:"required,min=3,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
}

type AssignRoleRequest struct {
	RoleID string `json:"role_id" binding:"required"`
}
