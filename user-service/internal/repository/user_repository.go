package repository

import (
	"database/sql"
	"github.com/Daty26/order-system/user-service/internal/model"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
}
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
func (pr *PostgresRepository) Create(user model.User) (model.User, error) {
	query := `Insert into users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := pr.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return model.User{}, err
	}
	return user, err

}
