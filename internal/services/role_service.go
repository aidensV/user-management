package services

import (
	"errors"

	"github.com/google/uuid"

	"user-management/internal/database"
	"user-management/internal/models"
	"user-management/internal/repositories"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repositories.NewRoleRepository(),
	}
}

// GetAllRoles returns all roles
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

// GetRoleByID returns a role by ID
func (s *RoleService) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return role, nil
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(name, displayName string, description *string) (*models.Role, error) {
	// Check if role already exists
	existing, err := s.roleRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("role already exists")
	}

	role := &models.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		IsSystem:    false,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(id uuid.UUID, displayName string, description *string, isActive bool) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// System roles cannot be deactivated
	if role.IsSystem && !isActive {
		return nil, errors.New("cannot deactivate system role")
	}

	role.DisplayName = displayName
	role.Description = description
	role.IsActive = isActive

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRole deletes a role (non-system only)
func (s *RoleService) DeleteRole(id uuid.UUID) error {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}
	if role.IsSystem {
		return errors.New("cannot delete system role")
	}
	return s.roleRepo.Delete(id)
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID, roleID, assignedBy uuid.UUID) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Check if already assigned
	var count int64
	db := database.GetDB()
	db.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
	if count > 0 {
		return errors.New("role already assigned to user")
	}

	userRole := &models.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: &assignedBy,
	}

	return db.Create(userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (s *RoleService) RemoveRoleFromUser(userID, roleID uuid.UUID) error {
	db := database.GetDB()
	return db.Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&models.UserRole{}).Error
}

// GetUserRoles returns all roles for a user
func (s *RoleService) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	db := database.GetDB()
	err := db.Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
