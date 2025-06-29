package models

import (
	"github.com/google/uuid"
	"time"
)

// Domain models for business logic
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // We don't expose in JSON
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Role         UserRole   `json:"role"`
	Status       UserStatus `json:"status"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type UserRole string

const (
	RoleGamer      UserRole = "gamer"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type UserStatus string

const (
	StatusActive    UserStatus = "active"
	StatusInactive  UserStatus = "inactive"
	StatusSuspended UserStatus = "suspended"
)

// Request/Response DTOs with validation
type CreateUserRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	FirstName string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string   `json:"last_name" validate:"required,min=2,max=100"`
	Role      UserRole `json:"role" validate:"required,oneof=buyer seller admin"`
	Phone     string   `json:"phone,omitempty" validate:"omitempty,min=10"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,min=10"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string     `json:"email" example:"user@example.com"`
	FirstName string     `json:"first_name" example:"John"`
	LastName  string     `json:"last_name" example:"Doe"`
	Role      UserRole   `json:"role" example:"buyer"`
	Status    UserStatus `json:"status" example:"active"`
	AvatarURL *string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	Phone     *string    `json:"phone,omitempty" example:"+1234567890"`
	CreatedAt time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func (r CreateUserRequest) GetSchema() interface{} {
	return r
}

func (r UpdateUserRequest) GetSchema() interface{} {
	return r
}

func (r LoginRequest) GetSchema() interface{} {
	return r
}
