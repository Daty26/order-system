package model

import "time"

type Roles string

var (
	UserRole  Roles = "USER"
	AdminRole Roles = "ADMIN"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      Roles     `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
