package models

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Resource    string    `gorm:"not null" json:"resource"`
	Action      string    `gorm:"not null" json:"action"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// GetPermissionKey returns "resource:action" format
func (p *Permission) GetKey() string {
	return p.Resource + ":" + p.Action
}

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName string    `gorm:"not null" json:"display_name"`
	Description *string   `json:"description,omitempty"`
	IsSystem    bool      `gorm:"default:false" json:"is_system"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type UserRole struct {
	UserID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	RoleID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"role_id"`
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty"`
	AssignedAt time.Time  `gorm:"default:now()" json:"assigned_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
