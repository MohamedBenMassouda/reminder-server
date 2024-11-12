package models

import "time"

type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// UserResponse for when you need to send user data with related entities
type UserResponse struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	Reminders []Reminder `json:"reminders,omitempty"`
}

type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
