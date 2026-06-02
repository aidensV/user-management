package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		db: database.DB,
	}
}

// Create creates a new menu
func (r *MenuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

// GetAll retrieves all menus
func (r *MenuRepository) GetAll() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Order("sort_order ASC").Find(&menus).Error
	return menus, err
}

// GetByID retrieves menu by ID
func (r *MenuRepository) GetByID(id uuid.UUID) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Where("id = ?", id).First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &menu, err
}

// Update updates a menu
func (r *MenuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

// Delete deletes a menu (cascade to children)
func (r *MenuRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Menu{}, "id = ?", id).Error
}

// GetMenuTree retrieves menu hierarchy
func (r *MenuRepository) GetMenuTree() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Where("is_active = true").Order("sort_order ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// Build tree
	menuMap := make(map[uuid.UUID]*models.Menu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	var roots []models.Menu
	for i := range menus {
		if menus[i].ParentID == nil {
			roots = append(roots, menus[i])
		} else {
			if parent, ok := menuMap[*menus[i].ParentID]; ok {
				parent.Children = append(parent.Children, menus[i])
			}
		}
	}

	return roots, nil
}

// GetMenusByRoleID retrieves menus accessible by a role
func (r *MenuRepository) GetMenusByRoleID(roleID uuid.UUID) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.
		Table("menus").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id = ? AND role_menus.can_view = true AND menus.is_active = true", roleID).
		Order("menus.sort_order ASC").
		Find(&menus).Error
	return menus, err
}

// AssignMenuToRole assigns menu access to a role
func (r *MenuRepository) AssignMenuToRole(roleID, menuID uuid.UUID, canView, canCreate, canEdit, canDelete bool) error {
	roleMenu := &models.RoleMenu{
		RoleID:    roleID,
		MenuID:    menuID,
		CanView:   canView,
		CanCreate: canCreate,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	}
	return r.db.Create(roleMenu).Error
}

// UpdateRoleMenu updates menu permissions for a role
func (r *MenuRepository) UpdateRoleMenu(roleID, menuID uuid.UUID, canView, canCreate, canEdit, canDelete bool) error {
	return r.db.Model(&models.RoleMenu{}).
		Where("role_id = ? AND menu_id = ?", roleID, menuID).
		Updates(map[string]interface{}{
			"can_view":   canView,
			"can_create": canCreate,
			"can_edit":   canEdit,
			"can_delete": canDelete,
		}).Error
}

// RemoveMenuFromRole removes menu access from a role
func (r *MenuRepository) RemoveMenuFromRole(roleID, menuID uuid.UUID) error {
	return r.db.Where("role_id = ? AND menu_id = ?", roleID, menuID).
		Delete(&models.RoleMenu{}).Error
}
