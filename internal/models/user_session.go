package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash        string    `gorm:"uniqueIndex;not null" json:"-"`
	RefreshTokenHash *string   `gorm:"uniqueIndex" json:"-"`
	IPAddress        *string   `json:"ip_address,omitempty"`
	UserAgent        *string   `json:"user_agent,omitempty"`
	DeviceName       *string   `json:"device_name,omitempty"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	ExpiresAt        time.Time `gorm:"not null;index" json:"expires_at"`
	LastActivityAt   time.Time `gorm:"default:now()" json:"last_activity_at"`
	CreatedAt        time.Time `json:"created_at"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
