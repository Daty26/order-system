package transport_http_dto

import (
	"time"

	"github.com/Daty26/order-system/user-service/internal/model"
)

type UserListItem struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Identifier string `json:"identifier"`
	Password string `json:"password"`
}
type LoginUserResponse struct{
	User model.UserSummary `json:"user"`
	Token string `json:"token"`
}
type CreateUserRequest struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}