package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-management/internal/services"
	"user-management/internal/utils"
)

type PermissionController struct {
	permissionService *services.PermissionService
	roleService       *services.RoleService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: services.NewPermissionService(),
		roleService:       services.NewRoleService(),
	}
}

// GetAllPermissions handles getting all permissions
// @Summary Get all permissions
// @Tags Permissions
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Router /api/v1/permissions [get]
func (c *PermissionController) GetAllPermissions(ctx *gin.Context) {
	permissions, err := c.permissionService.GetAllPermissions()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_PERMISSIONS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, permissions, "Permissions retrieved")
}

// CreatePermission handles creating a new permission
// @Summary Create permission
// @Tags Permissions
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param request body CreatePermissionRequest true "Permission data"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/permissions [post]
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	permission, err := c.permissionService.CreatePermission(
		req.Resource,
		req.Action,
		req.Name,
		req.Description,
	)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "CREATE_PERMISSION_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, permission, "Permission created successfully")
}

// GetPermissionsByRole handles getting permissions for a role
// @Summary Get permissions by role
// @Tags Permissions
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/permissions [get]
func (c *PermissionController) GetPermissionsByRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	permissions, err := c.permissionService.GetPermissionsByRole(roleID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_PERMISSIONS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, permissions, "Permissions retrieved")
}

// AssignPermissionToRole handles assigning a permission to a role
// @Summary Assign permission to role
// @Tags Permissions
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Param permissionId path string true "Permission ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/permissions/{permissionId} [post]
func (c *PermissionController) AssignPermissionToRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	permissionIDParam := ctx.Param("permissionId")
	permissionID, err := uuid.Parse(permissionIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid permission ID", "INVALID_PERMISSION_ID")
		return
	}

	if err := c.permissionService.AssignPermissionToRole(roleID, permissionID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "ASSIGN_PERMISSION_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Permission assigned successfully")
}

// RemovePermissionFromRole handles removing a permission from a role
// @Summary Remove permission from role
// @Tags Permissions
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Param permissionId path string true "Permission ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/permissions/{permissionId} [delete]
func (c *PermissionController) RemovePermissionFromRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	permissionIDParam := ctx.Param("permissionId")
	permissionID, err := uuid.Parse(permissionIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid permission ID", "INVALID_PERMISSION_ID")
		return
	}

	if err := c.permissionService.RemovePermissionFromRole(roleID, permissionID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "REMOVE_PERMISSION_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Permission removed successfully")
}

// CreatePermissionRequest DTO
type CreatePermissionRequest struct {
	Resource    string  `json:"resource" binding:"required,min=2,max=100"`
	Action      string  `json:"action" binding:"required,min=2,max=50"`
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description,omitempty"`
}
