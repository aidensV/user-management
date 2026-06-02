package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type LoginHistoryRepository struct {
	db *gorm.DB
}

func NewLoginHistoryRepository() *LoginHistoryRepository {
	return &LoginHistoryRepository{
		db: database.DB,
	}
}

// Create creates a new login history record
func (r *LoginHistoryRepository) Create(history *models.LoginHistory) error {
	return r.db.Create(history).Error
}

// GetByUserID retrieves login history for a user
func (r *LoginHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]models.LoginHistory, error) {
	var histories []models.LoginHistory
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&histories).Error
	return histories, err
}

// GetAll retrieves all login history with pagination
func (r *LoginHistoryRepository) GetAll(page, pageSize int, search string) ([]models.LoginHistory, int64, error) {
	var histories []models.LoginHistory
	var total int64

	query := r.db.Model(&models.LoginHistory{})

	if search != "" {
		query = query.Where("email ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// CleanOldRecords removes records older than days
func (r *LoginHistoryRepository) CleanOldRecords(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&models.LoginHistory{}).Error
}
