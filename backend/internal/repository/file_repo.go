package repository

import (
	"astro-planner/backend/internal/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IFileRepository interface {
}

type FileRepo struct {
	db *pgxpool.Pool
}

func NewFileRepo(db *pgxpool.Pool) *FileRepo {
	return &FileRepo{db: db}
}

func (r *FileRepo) Save(ctx context.Context, f *entities.File) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO files (user_id, object_key, original_name, content_type, size)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING file_id, created_at
	`, f.UserID, f.ObjectKey, f.OriginalName, f.ContentType, f.Size).Scan(
		&f.ID, &f.CreatedAt,
	)
}
