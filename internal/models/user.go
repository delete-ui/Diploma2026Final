package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                    uuid.UUID  `json:"id"`
	Email                 string     `json:"email"`
	Username              string     `json:"username"`
	PasswordHash          string     `json:"-"`
	Role                  string     `json:"role"`
	IsVerified            bool       `json:"is_verified"`
	VerificationCode      *string    `json:"-"`
	VerificationExpiresAt *time.Time `json:"-"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	User        User   `json:"user"`
}
