package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail gets user by email (including soft deleted)
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByID gets user by ID (only active, not soft deleted)
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByIDUnscoped gets user by ID including soft deleted
func (r *UserRepository) GetByIDUnscoped(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all active users with pagination
func (r *UserRepository) GetAll(page, pageSize int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{}).Where("deleted_at IS NULL")

	if search != "" {
		query = query.Where("email ILIKE ? OR name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// SoftDelete soft deletes a user
func (r *UserRepository) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"deleted_by": deletedBy,
			"is_active":  false,
		}).Error
}

// UpdateLastLogin updates user's last login info
func (r *UserRepository) UpdateLastLogin(id uuid.UUID, ip string) error {
	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		}).Error
}

// IncrementFailedLoginAttempts increments failed login attempts
func (r *UserRepository) IncrementFailedLoginAttempts(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).Error
}

// ResetFailedLoginAttempts resets failed login attempts
func (r *UserRepository) ResetFailedLoginAttempts(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("failed_login_attempts", 0).Error
}

// LockAccount locks a user account
func (r *UserRepository) LockAccount(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_locked": true,
			"is_active": false,
		}).Error
}

// GetByUsername gets user by username (including soft deleted)
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
