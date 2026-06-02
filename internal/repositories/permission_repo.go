package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{
		db: database.DB,
	}
}

// Create creates a new permission
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// GetAll retrieves all permissions
func (r *PermissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Order("resource, action").Find(&permissions).Error
	return permissions, err
}

// GetByID retrieves permission by ID
func (r *PermissionRepository) GetByID(id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &permission, err
}

// GetByResourceAndAction retrieves permission by resource and action
func (r *PermissionRepository) GetByResourceAndAction(resource, action string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("resource = ? AND action = ?", resource, action).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &permission, err
}

// GetPermissionsByRoleID retrieves all permissions for a role
func (r *PermissionRepository) GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// GetPermissionsByUserID retrieves all permissions for a user (from their roles)
func (r *PermissionRepository) GetPermissionsByUserID(userID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct().
		Find(&permissions).Error
	return permissions, err
}
