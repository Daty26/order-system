package repository

import (
	"database/sql"
	"errors"
	"github.com/Daty26/order-system/user-service/internal/model"
	"github.com/lib/pq"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
	GetByUsernameOrEmail(identifier string) (model.User, error)
}
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (pr *PostgresRepository) Create(user model.User) (model.User, error) {
	query := `Insert into users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := pr.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			if pgError.Constraint == "users_name_key" {
				return model.User{}, errors.New("username already exists")
			}
			if pgError.Constraint == "users_email_key" {
				return model.User{}, errors.New("email already exists")
			}
			return model.User{}, errors.New("duplicate value for a unique field")
		}
		return model.User{}, err
	}
	return user, err
}
func (pr *PostgresRepository) GetByUsernameOrEmail(identifier string) (model.User, error) {
	var user model.User
	query := `Select id, username, email, password, role, created_at from users where username=$1 or email =$1`
	if err := pr.db.QueryRow(query, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
		return model.User{}, err
	}
	return user, nil
}
