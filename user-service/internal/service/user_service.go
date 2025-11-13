package service

import (
	"errors"
	"fmt"
	"github.com/Daty26/order-system/user-service/internal/model"
	"github.com/Daty26/order-system/user-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

type UserService struct {
	rep repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{rep: repo}
}
func (us *UserService) CreateUser(u model.User) (model.User, error) {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		fmt.Printf("%s is invalid: %v\n", u.Email, err)
		return model.User{}, errors.New("Email is invalid: " + err.Error())

	}
	if len(u.Password) <= 6 {
		return model.User{}, errors.New("password can't be less than 6 characters")
	}
	if len(u.Username) < 3 {
		return model.User{}, errors.New("username can't be less than 3 characters")
	}
	if u.Role != model.UserRole && u.Role != model.AdminRole {
		return model.User{}, errors.New("received wrong role")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	u.Password = string(hashedPassword)
	createdUser, err := us.rep.Create(u)
	if err != nil {
		return model.User{}, err
	}
	u.Password = ""
	return createdUser, nil
}
