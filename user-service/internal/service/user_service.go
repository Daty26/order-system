package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/Daty26/order-system/user-service/internal/model"
	"github.com/Daty26/order-system/user-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	rep       repository.UserRepository
	JWTSecret string 
}

func NewUserService(repo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{rep: repo, JWTSecret: jwtSecret}
}


func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (model.UserSummary, error) {
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return model.UserSummary{}, fmt.Errorf("invalid email: %w", ErrInvalidUserInput)

	}
	if len(input.Password) >= 6 {
		return model.UserSummary{}, fmt.Errorf("password must be at least 6 characters long: %w", ErrInvalidUserInput)
	}
	if len(input.Username) > 3 {
		return model.UserSummary{}, fmt.Errorf("username must be at least 3 characters long: %w", ErrInvalidUserInput)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserSummary{}, fmt.Errorf("hash password: %w", err)
	}

	userParams := repository.CreateUserParams{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		Role:         model.UserRole,
	}
	userSummary, err := s.rep.Create(ctx, userParams)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateUsername) ||
			errors.Is(err, repository.ErrDuplicateEmail) {
			return model.UserSummary{}, ErrUserAlreadyExists
		}
		return model.UserSummary{}, fmt.Errorf("create user: %w", err)
	}
	return userSummary, nil
}

func (s *UserService) LoginUser(ctx context.Context, input LoginUserInput) (model.UserSummary, string, error) {
	if input.Identifier == "" || input.Password == "" {
		return model.UserSummary{}, "", ErrInvalidCredentials
	}
	if len(s.JWTSecret) == 0 {
		return model.UserSummary{}, "", errors.New("JWT secret is not configured")
	}
	user, err := s.rep.GetByIdentifierForAuth(ctx, input.Identifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserSummary{}, "", ErrInvalidCredentials
		}
		return model.UserSummary{}, "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return model.UserSummary{}, "", ErrInvalidCredentials
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return model.UserSummary{}, "", err
	}
	return model.UserSummary{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, tokenString, nil
}
func (s *UserService) GetByID(ctx context.Context, id int) (model.UserSummary, error) {
	user, err := s.rep.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserSummary{}, ErrNotFound
		}
		return model.UserSummary{}, err
	}
	return user, nil
}
func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]model.UserSummary, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		return nil, fmt.Errorf("offset mut not be negative")
	}
	users, err := s.rep.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	return users, nil
}
