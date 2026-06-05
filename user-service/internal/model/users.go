package model

import "time"

type Role string

var (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	PasswordHash  string    `json:"-"`
	Role      Role     `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSummary struct {
	ID        int
	Username  string
	Email     string
	Role      Role
	CreatedAt time.Time
}
