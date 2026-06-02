package models

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ParentID  *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Name      string     `gorm:"not null" json:"name"`
	Icon      *string    `json:"icon,omitempty"`
	Path      *string    `json:"path,omitempty"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relations
	Children []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Menu) TableName() string {
	return "menus"
}

type RoleMenu struct {
	RoleID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	MenuID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"menu_id"`
	CanView   bool      `gorm:"default:true" json:"can_view"`
	CanCreate bool      `gorm:"default:false" json:"can_create"`
	CanEdit   bool      `gorm:"default:false" json:"can_edit"`
	CanDelete bool      `gorm:"default:false" json:"can_delete"`
	CreatedAt time.Time `json:"created_at"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}

type LoginHistory struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Email         string     `gorm:"not null" json:"email"`
	Status        string     `gorm:"not null;check:status IN ('success','failed','locked')" json:"status"`
	IPAddress     *string    `json:"ip_address,omitempty"`
	UserAgent     *string    `json:"user_agent,omitempty"`
	FailureReason *string    `json:"failure_reason,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (LoginHistory) TableName() string {
	return "login_history"
}

type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `json:"created_at"`

	// Relation
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
