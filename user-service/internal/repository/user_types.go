package repository

import "github.com/Daty26/order-system/user-service/internal/model"

type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
	Role         model.Role
}
