package repository

import (
	"astro-planner/backend/internal/entities"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	UpdateCoordinates(ctx context.Context, userId string, latitude, lotitude float64) error
}

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
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

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.QueryRow(ctx, `
	SELECT user_id, name, email, password_hash, role, is_blocked, register_date, latitude, longitude
	FROM users
	WHERE user_id = $1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsBlocked,
		&user.RegisterDate,
		&user.Latitude,
		&user.Longitude,
	)
	if err != nil {
		log.Println("ОШИБКА В БАЗЕ (GetByID):", err)
		return nil, err
	}
	return user, err
}

func (r *PostgresUserRepository) UpdateCoordinates(ctx context.Context, userId string, latitude, longitude float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users  
		SET latitude = $1, longitude = $2 
		WHERE user_id = $3`, latitude, longitude, userId)
	return err
}

func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, userId, newPasswordHash string) error {
	_, err := r.db.Exec(ctx, `
	UPDATE users
	SET password_hash = $1
	WHERE user_id = $2`, newPasswordHash, userId)
	return err
}

func (r *PostgresUserRepository) UpdateRole(ctx context.Context, userId, newRole string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET role = $1::role_enum WHERE user_id = $2`, newRole, userId)
	return err
}

func (r *PostgresUserRepository) SetBlockUser(ctx context.Context, userId string, isBlocked bool) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET is_blocked = $1 WHERE user_id = $2`, isBlocked, userId)
	return err
}
