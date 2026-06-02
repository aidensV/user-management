package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		db: database.DB,
	}
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// GetAll retrieves all roles
func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Order("name").Find(&roles).Error
	return roles, err
}

// GetByID retrieves role by ID
func (r *RoleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// GetByName retrieves role by name
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role (only non-system roles)
func (r *RoleRepository) Delete(id uuid.UUID) error {
	// Check if role is system
	var role models.Role
	if err := r.db.Where("id = ? AND is_system = ?", id, false).First(&role).Error; err != nil {
		return errors.New("cannot delete system role")
	}
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}

// AssignPermission assigns a permission to a role
func (r *RoleRepository) AssignPermission(roleID, permissionID uuid.UUID) error {
	rolePermission := &models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(rolePermission).Error
}

// RemovePermission removes a permission from a role
func (r *RoleRepository) RemovePermission(roleID, permissionID uuid.UUID) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}

// GetRolePermissions retrieves all permission IDs for a role
func (r *RoleRepository) GetRolePermissions(roleID uuid.UUID) ([]uuid.UUID, error) {
	var permissionIDs []uuid.UUID
	err := r.db.Model(&models.RolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	return permissionIDs, err
}
