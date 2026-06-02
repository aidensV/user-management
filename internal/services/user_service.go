package services

import (
	"errors"

	"github.com/google/uuid"

	"user-management/internal/models"
	"user-management/internal/repositories"
	"user-management/internal/utils"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

// CreateUser creates a new user (admin only)
func (s *UserService) CreateUser(req *models.CreateUserRequest, createdBy uuid.UUID) (*models.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if existingUser != nil && existingUser.DeletedAt == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Email:      req.Email,
		Username:   req.Username,
		Password:   hashedPassword,
		Name:       req.Name,
		NIK:        req.NIK,
		Department: req.Department,
		Position:   req.Position,
		Phone:      req.Phone,
		Role:       req.Role,
		IsActive:   true,
		CreatedBy:  &createdBy,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(page, pageSize int, search string) ([]models.UserResponse, int64, error) {
	users, total, err := s.userRepo.GetAll(page, pageSize, search)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve users")
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, total, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest, updatedBy uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("internal server error")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Username != nil {
		user.Username = req.Username
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.NIK != nil {
		user.NIK = req.NIK
	}
	if req.Department != nil {
		user.Department = req.Department
	}
	if req.Position != nil {
		user.Position = req.Position
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedBy = &updatedBy

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id uuid.UUID, deletedBy uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("internal server error")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Prevent deleting yourself
	if user.ID == deletedBy {
		return errors.New("cannot delete your own account")
	}

	// Prevent deleting the last admin (optional - can be enhanced)
	if user.Role == "admin" {
		// Check if there's at least one other admin
		admins, _, err := s.userRepo.GetAll(1, 100, "")
		if err != nil {
			return errors.New("failed to check admin count")
		}
		adminCount := 0
		for _, a := range admins {
			if a.Role == "admin" && a.ID != id && a.DeletedAt == nil {
				adminCount++
			}
		}
		if adminCount == 0 {
			return errors.New("cannot delete the last admin user")
		}
	}

	return s.userRepo.SoftDelete(id, deletedBy)
}

// GetUserWithRoles retrieves user with their roles
func (s *UserService) GetUserWithRoles(id uuid.UUID) (*models.UserResponse, []models.Role, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	// Get user's roles
	roleService := NewRoleService()
	roles, err := roleService.GetUserRoles(id)
	if err != nil {
		return nil, nil, err
	}

	response := user.ToResponse()
	return &response, roles, nil
}

// AssignRoleToUser assigns a role to a user
func (s *UserService) AssignRoleToUser(userID, roleID, assignedBy uuid.UUID) error {
	roleService := NewRoleService()
	return roleService.AssignRoleToUser(userID, roleID, assignedBy)
}

// RemoveRoleFromUser removes a role from a user
func (s *UserService) RemoveRoleFromUser(userID, roleID uuid.UUID) error {
	roleService := NewRoleService()
	return roleService.RemoveRoleFromUser(userID, roleID)
}
