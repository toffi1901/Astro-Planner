package repository

import (
	"astro-planner/backend/internal/entities"
	"context"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewpostgresUsereposotory(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.QueryRow(ctx, `
		SELECT user_id, name, email, password_hash, role, is_blocked, register_date
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsBlocked,
		&user.RegisterDate,
	)
	return user, err
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.QueryRow(ctx, `INSERT INTO users (email, name, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, register_date`, user.Email, user.Name, user.PasswordHash, user.Role).Scan(
		&user.ID,
		&user.RegisterDate,
	)
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE user_id=$1`, id)
	return err
}
