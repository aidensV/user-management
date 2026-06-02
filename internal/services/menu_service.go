package services

import (
	"errors"

	"github.com/google/uuid"

	"user-management/internal/database"
	"user-management/internal/models"
	"user-management/internal/repositories"
)

type MenuService struct {
	menuRepo *repositories.MenuRepository
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuRepo: repositories.NewMenuRepository(),
	}
}

// GetAllMenus returns all menus
func (s *MenuService) GetAllMenus() ([]models.Menu, error) {
	return s.menuRepo.GetAll()
}

// GetMenuTree returns menu hierarchy
func (s *MenuService) GetMenuTree() ([]models.Menu, error) {
	return s.menuRepo.GetMenuTree()
}

// GetMenusByRole returns menus accessible by a role
func (s *MenuService) GetMenusByRole(roleID uuid.UUID) ([]models.Menu, error) {
	role, err := repositories.NewRoleRepository().GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return s.menuRepo.GetMenusByRoleID(roleID)
}

// GetMenusByUser returns menus accessible by a user (from their roles)
func (s *MenuService) GetMenusByUser(userID uuid.UUID) ([]models.Menu, error) {
	// Get user's roles
	var roles []models.Role
	db := database.GetDB()
	err := db.Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	// Collect menus from all roles
	menuMap := make(map[uuid.UUID]models.Menu)
	for _, role := range roles {
		menus, err := s.menuRepo.GetMenusByRoleID(role.ID)
		if err != nil {
			continue
		}
		for _, menu := range menus {
			menuMap[menu.ID] = menu
		}
	}

	// Convert map to slice
	result := make([]models.Menu, 0, len(menuMap))
	for _, menu := range menuMap {
		result = append(result, menu)
	}

	return result, nil
}

// CreateMenu creates a new menu
func (s *MenuService) CreateMenu(name string, parentID *uuid.UUID, icon, path *string, sortOrder int) (*models.Menu, error) {
	menu := &models.Menu{
		Name:      name,
		ParentID:  parentID,
		Icon:      icon,
		Path:      path,
		SortOrder: sortOrder,
		IsActive:  true,
	}

	if err := s.menuRepo.Create(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

// UpdateMenu updates an existing menu
func (s *MenuService) UpdateMenu(id uuid.UUID, name string, parentID *uuid.UUID, icon, path *string, sortOrder int, isActive bool) (*models.Menu, error) {
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if menu == nil {
		return nil, errors.New("menu not found")
	}

	menu.Name = name
	menu.ParentID = parentID
	menu.Icon = icon
	menu.Path = path
	menu.SortOrder = sortOrder
	menu.IsActive = isActive

	if err := s.menuRepo.Update(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

// DeleteMenu deletes a menu
func (s *MenuService) DeleteMenu(id uuid.UUID) error {
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return err
	}
	if menu == nil {
		return errors.New("menu not found")
	}
	return s.menuRepo.Delete(id)
}

// AssignMenuToRole assigns menu access to a role
func (s *MenuService) AssignMenuToRole(roleID, menuID uuid.UUID, canView, canCreate, canEdit, canDelete bool) error {
	// Check if already exists
	var existing models.RoleMenu
	db := database.GetDB()
	err := db.Where("role_id = ? AND menu_id = ?", roleID, menuID).First(&existing).Error
	if err == nil {
		return s.menuRepo.UpdateRoleMenu(roleID, menuID, canView, canCreate, canEdit, canDelete)
	}
	return s.menuRepo.AssignMenuToRole(roleID, menuID, canView, canCreate, canEdit, canDelete)
}

// RemoveMenuFromRole removes menu access from a role
func (s *MenuService) RemoveMenuFromRole(roleID, menuID uuid.UUID) error {
	return s.menuRepo.RemoveMenuFromRole(roleID, menuID)
}
