package models

import "time"

type User struct {
	ID           int    `db:"id"            json:"id"`
	Name         string `db:"name"          json:"name"`
	Email        string `db:"email"         json:"email"`
	Avatar       string `db:"avatar"        json:"avatar"`
	PasswordHash string `db:"password_hash" json:"-"`
}

type CreateUserInput struct {
	Name     string `json:"name"     binding:"required,min=2,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type UpdateUserInput struct {
	Name  string `json:"name"  binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PaginatedResponse struct {
	Data  any `json:"data"`
	Page  int `json:"page"`
	Total int `json:"total"`
	Limit int `json:"limit"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
}
