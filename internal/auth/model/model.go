package model

import (
	"time"
)

// DB Models
type User struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"-"` // hide password from JSON
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// DTOs
type AuthRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"-" validate:"required,min=8,max=32"`
}

type AuthResponse struct {
	User  User   `json:"user" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyTokenResponse struct {
	User User `json:"user" validate:"required"`
}