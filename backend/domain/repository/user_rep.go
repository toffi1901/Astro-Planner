package repository

import (
	"context"
	"go/constant"
)

type UserRepository interface{
	GetById(ctx context.Context, id uuid) (User, error)
	Create(ctx context.Context, user User) error
	Delete(ctx context.Context, id uuid) error
}

type PostgresUserrepository struct {
	db *pgxpool.Pool
}

func NewpostgresUsereposotory(db *pgxpool.pool) *PostgresUserrepository{
	return &PostgresUserrepository{
		db: db,
	}
}

func (r *PostgreUserrepository) GetById(ctx context.Context, id uuid) (User, error){
	var user UserRepositoryerr := r.db.QueryRow()
}