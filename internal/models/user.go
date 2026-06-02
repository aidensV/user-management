package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email      string    `gorm:"uniqueIndex;not null" json:"email"`
	Username   *string   `gorm:"uniqueIndex" json:"username,omitempty"`
	Password   string    `gorm:"not null" json:"-"` // "-" means exclude from JSON
	Name       string    `gorm:"not null" json:"name"`
	NIK        *string   `gorm:"uniqueIndex" json:"nik,omitempty"`
	Department *string   `json:"department,omitempty"`
	Position   *string   `json:"position,omitempty"`
	Phone      *string   `json:"phone,omitempty"`

	// Role (temporary for P0)
	Role string `gorm:"not null;check:role IN ('admin','supervisor','warehouse_operator','viewer')" json:"role"`

	// Status
	IsActive            bool       `gorm:"default:true" json:"is_active"`
	IsLocked            bool       `gorm:"default:false" json:"is_locked"`
	FailedLoginAttempts int        `gorm:"default:0" json:"-"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP         *string    `json:"last_login_ip,omitempty"`

	// Soft delete
	DeletedAt *time.Time `gorm:"index" json:"-"`
	DeletedBy *uuid.UUID `json:"-"`

	// Audit
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedBy *uuid.UUID `json:"-"`
	UpdatedBy *uuid.UUID `json:"-"`
}

func (User) TableName() string {
	return "users"
}

// UserResponse is the struct returned to client (excludes sensitive fields)
type UserResponse struct {
	ID          uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	Username    *string    `json:"username,omitempty"`
	Name        string     `json:"name"`
	NIK         *string    `json:"nik,omitempty"`
	Department  *string    `json:"department,omitempty"`
	Position    *string    `json:"position,omitempty"`
	Phone       *string    `json:"phone,omitempty"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		Name:        u.Name,
		NIK:         u.NIK,
		Department:  u.Department,
		Position:    u.Position,
		Phone:       u.Phone,
		Role:        u.Role,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// CreateUserRequest (Admin only)
type CreateUserRequest struct {
	Email      string  `json:"email" binding:"required,email"`
	Username   *string `json:"username"`
	Password   string  `json:"password" binding:"required,min=6"`
	Name       string  `json:"name" binding:"required"`
	NIK        *string `json:"nik"`
	Department *string `json:"department"`
	Position   *string `json:"position"`
	Phone      *string `json:"phone"`
	Role       string  `json:"role" binding:"required,oneof=admin supervisor warehouse_operator viewer"`
}

// UpdateUserRequest
type UpdateUserRequest struct {
	Username   *string `json:"username"`
	Name       *string `json:"name"`
	NIK        *string `json:"nik"`
	Department *string `json:"department"`
	Position   *string `json:"position"`
	Phone      *string `json:"phone"`
	Role       *string `json:"role" binding:"omitempty,oneof=admin supervisor warehouse_operator viewer"`
	IsActive   *bool   `json:"is_active"`
}
