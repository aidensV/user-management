package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository() *PasswordResetRepository {
	return &PasswordResetRepository{
		db: database.DB,
	}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// GetByToken retrieves reset token by token string
func (r *PasswordResetRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ? AND is_used = false AND expires_at > ?", token, time.Now()).
		First(&resetToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &resetToken, err
}

// MarkAsUsed marks a token as used
func (r *PasswordResetRepository) MarkAsUsed(id uuid.UUID) error {
	return r.db.Model(&models.PasswordResetToken{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

// CleanExpiredTokens removes expired tokens
func (r *PasswordResetRepository) CleanExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
}
