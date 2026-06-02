package services

import (
	"errors"

	"github.com/google/uuid"

	"user-management/internal/models"
	"user-management/internal/repositories"
)

type PermissionService struct {
	permissionRepo *repositories.PermissionRepository
	roleRepo       *repositories.RoleRepository
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		permissionRepo: repositories.NewPermissionRepository(),
		roleRepo:       repositories.NewRoleRepository(),
	}
}

// GetAllPermissions returns all permissions
func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	return s.permissionRepo.GetAll()
}

// GetPermissionsByRole returns permissions for a specific role
func (s *PermissionService) GetPermissionsByRole(roleID uuid.UUID) ([]models.Permission, error) {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return s.permissionRepo.GetPermissionsByRoleID(roleID)
}

// GetPermissionsByUser returns permissions for a specific user
func (s *PermissionService) GetPermissionsByUser(userID uuid.UUID) ([]models.Permission, error) {
	return s.permissionRepo.GetPermissionsByUserID(userID)
}

// AssignPermissionToRole assigns a permission to a role
func (s *PermissionService) AssignPermissionToRole(roleID, permissionID uuid.UUID) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	permission, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		return err
	}
	if permission == nil {
		return errors.New("permission not found")
	}

	return s.roleRepo.AssignPermission(roleID, permissionID)
}

// RemovePermissionFromRole removes a permission from a role
func (s *PermissionService) RemovePermissionFromRole(roleID, permissionID uuid.UUID) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.roleRepo.RemovePermission(roleID, permissionID)
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(resource, action, name string, description *string) (*models.Permission, error) {
	// Check if permission already exists
	existing, err := s.permissionRepo.GetByResourceAndAction(resource, action)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("permission already exists")
	}

	permission := &models.Permission{
		Resource:    resource,
		Action:      action,
		Name:        name,
		Description: description,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

// UserHasPermission checks if a user has a specific permission
func (s *PermissionService) UserHasPermission(userID uuid.UUID, resource, action string) (bool, error) {
	permissions, err := s.permissionRepo.GetPermissionsByUserID(userID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p.Resource == resource && p.Action == action {
			return true, nil
		}
	}
	return false, nil
}
