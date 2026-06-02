package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-management/internal/services"
	"user-management/internal/utils"
)

type MenuController struct {
	menuService *services.MenuService
}

func NewMenuController() *MenuController {
	return &MenuController{
		menuService: services.NewMenuService(),
	}
}

// GetAllMenus handles getting all menus
// @Summary Get all menus
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Router /api/v1/menus [get]
func (c *MenuController) GetAllMenus(ctx *gin.Context) {
	menus, err := c.menuService.GetAllMenus()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_MENUS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, menus, "Menus retrieved")
}

// GetMenuTree handles getting menu hierarchy
// @Summary Get menu tree
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Router /api/v1/menus/tree [get]
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	menus, err := c.menuService.GetMenuTree()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_MENU_TREE_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, menus, "Menu tree retrieved")
}

// GetUserMenus handles getting menus for current user
// @Summary Get user menus
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Success 200 {object} utils.Response
// @Router /api/v1/menus/user [get]
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
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

	menus, err := c.menuService.GetMenusByUser(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_USER_MENUS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, menus, "User menus retrieved")
}

// GetMenusByRole handles getting menus for a specific role
// @Summary Get menus by role
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/menus [get]
func (c *MenuController) GetMenusByRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	menus, err := c.menuService.GetMenusByRole(roleID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), "GET_ROLE_MENUS_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, menus, "Role menus retrieved")
}

// CreateMenu handles creating a new menu
// @Summary Create menu
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param request body CreateMenuRequest true "Menu data"
// @Success 201 {object} utils.Response
// @Router /api/v1/menus [post]
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var req CreateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	menu, err := c.menuService.CreateMenu(req.Name, req.ParentID, req.Icon, req.Path, req.SortOrder)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "CREATE_MENU_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, menu, "Menu created successfully")
}

// UpdateMenu handles updating a menu
// @Summary Update menu
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param id path string true "Menu ID"
// @Param request body UpdateMenuRequest true "Menu data"
// @Success 200 {object} utils.Response
// @Router /api/v1/menus/{id} [put]
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	idParam := ctx.Param("id")
	menuID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid menu ID", "INVALID_ID")
		return
	}

	var req UpdateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	menu, err := c.menuService.UpdateMenu(menuID, req.Name, req.ParentID, req.Icon, req.Path, req.SortOrder, req.IsActive)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "UPDATE_MENU_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, menu, "Menu updated successfully")
}

// DeleteMenu handles deleting a menu
// @Summary Delete menu
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param id path string true "Menu ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/menus/{id} [delete]
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	idParam := ctx.Param("id")
	menuID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid menu ID", "INVALID_ID")
		return
	}

	if err := c.menuService.DeleteMenu(menuID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "DELETE_MENU_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Menu deleted successfully")
}

// AssignMenuToRole handles assigning menu to role
// @Summary Assign menu to role
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Param menuId path string true "Menu ID"
// @Param request body AssignMenuRequest true "Menu permissions"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/menus/{menuId} [post]
func (c *MenuController) AssignMenuToRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	menuIDParam := ctx.Param("menuId")
	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid menu ID", "INVALID_MENU_ID")
		return
	}

	var req AssignMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request: "+err.Error(), "INVALID_REQUEST")
		return
	}

	err = c.menuService.AssignMenuToRole(roleID, menuID, req.CanView, req.CanCreate, req.CanEdit, req.CanDelete)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "ASSIGN_MENU_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Menu assigned to role successfully")
}

// RemoveMenuFromRole handles removing menu from role
// @Summary Remove menu from role
// @Tags Menus
// @Security ApiKeyAuth
// @param Authorization header string true "Bearer token"
// @Param roleId path string true "Role ID"
// @Param menuId path string true "Menu ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/roles/{roleId}/menus/{menuId} [delete]
func (c *MenuController) RemoveMenuFromRole(ctx *gin.Context) {
	roleIDParam := ctx.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid role ID", "INVALID_ROLE_ID")
		return
	}

	menuIDParam := ctx.Param("menuId")
	menuID, err := uuid.Parse(menuIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid menu ID", "INVALID_MENU_ID")
		return
	}

	if err := c.menuService.RemoveMenuFromRole(roleID, menuID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), "REMOVE_MENU_FAILED")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, nil, "Menu removed from role successfully")
}

// Request DTOs
type CreateMenuRequest struct {
	Name      string     `json:"name" binding:"required,min=2,max=100"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Icon      *string    `json:"icon,omitempty"`
	Path      *string    `json:"path,omitempty"`
	SortOrder int        `json:"sort_order"`
}

type UpdateMenuRequest struct {
	Name      string     `json:"name" binding:"required,min=2,max=100"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Icon      *string    `json:"icon,omitempty"`
	Path      *string    `json:"path,omitempty"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
}

type AssignMenuRequest struct {
	CanView   bool `json:"can_view"`
	CanCreate bool `json:"can_create"`
	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
}
