package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daty26/order-system/user-service/internal/model"
	"github.com/lib/pq"
)

type UserRepository interface {
	Create(ctx context.Context, userParams CreateUserParams) (model.UserSummary, error)
	GetByIdentifierForAuth(ctx context.Context,identifier string) (model.User, error)
	GetByID(ctx context.Context, id int) (model.UserSummary, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.UserSummary, error)
}
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (pr *PostgresRepository) Create(ctx context.Context, input CreateUserParams) (model.UserSummary, error) {
	var user model.UserSummary
	query := `
	Insert into users (username, email, password, role) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id, username, email, role, created_at`
	err := pr.db.QueryRowContext(
		ctx,
		query,
		input.Username, 
		input.Email, 
		input.PasswordHash, 
		input.Role,
		).Scan(&user.ID,&user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			switch pgError.Constraint{
			case "users_usersname_key":
				return model.UserSummary{}, ErrDuplicateUsername
			case "users_email_key":
				return model.UserSummary{}, ErrDuplicateEmail
			default:
				return model.UserSummary{}, fmt.Errorf("%w: %s", ErrDuplicateUser, pgError.Constraint)
			}
		}
		return model.UserSummary{}, fmt.Errorf("insert user: %w", err)
	}
	return user, nil
}
func (pr *PostgresRepository) GetByIdentifierForAuth(ctx context.Context,identifier string) (model.User, error) {
	var user model.User
	query := `Select id, username, email, password, role, created_at from users where username=$1 or email =$1`
	if err := pr.db.QueryRowContext(ctx, query, identifier).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		return model.User{}, err
	}
	return user, nil
}
func (pr *PostgresRepository) GetByID(ctx context.Context, id int) (model.UserSummary, error) {
	var user model.UserSummary
	query := `Select id, username, email, role, created_at from users where id=$1`
	if err := pr.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email,  &user.Role, &user.CreatedAt); err != nil {
		return model.UserSummary{}, err
	}
	return user, nil
}
func (r *PostgresRepository) GetAll(ctx context.Context, limit, offset int) ([]model.UserSummary, error) {
	query := `
	SELECT id, username, email, role, created_at FROM users
	ORDER BY id
	LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return []model.UserSummary{}, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()
	users := make([]model.UserSummary, 0)
	for rows.Next() {
		var user model.UserSummary
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
			return []model.UserSummary{}, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}
	return users, nil
}
