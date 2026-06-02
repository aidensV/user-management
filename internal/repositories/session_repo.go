package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"user-management/internal/database"
	"user-management/internal/models"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		db: database.DB,
	}
}

// Create creates a new user session
func (r *SessionRepository) Create(session *models.UserSession) error {
	return r.db.Create(session).Error
}

// GetByTokenHash gets session by token hash
func (r *SessionRepository) GetByTokenHash(tokenHash string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.Where("token_hash = ? AND is_active = true", tokenHash).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

// DeleteByTokenHash deletes session by token hash (logout)
func (r *SessionRepository) DeleteByTokenHash(tokenHash string) error {
	return r.db.Where("token_hash = ?", tokenHash).Delete(&models.UserSession{}).Error
}

// DeleteAllByUserID deletes all sessions for a user (force logout)
func (r *SessionRepository) DeleteAllByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}

// GetActiveSessionsByUserID gets all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(userID uuid.UUID) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := r.db.Where("user_id = ? AND is_active = true AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

// CleanExpiredSessions removes expired sessions
func (r *SessionRepository) CleanExpiredSessions() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{}).Error
}

// UpdateLastActivity updates last activity timestamp
func (r *SessionRepository) UpdateLastActivity(tokenHash string) error {
	return r.db.Model(&models.UserSession{}).
		Where("token_hash = ?", tokenHash).
		Update("last_activity_at", time.Now()).Error
}
