package repository

import (
	"astro-planner/backend/internal/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CelestialRepo struct{ db *pgxpool.Pool }

func NewCelestialRepo(db *pgxpool.Pool) *CelestialRepo { return &CelestialRepo{db: db} }

func (r *CelestialRepo) List(ctx context.Context) ([]entities.CelestialObject, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.object_id, c.name, c.object_type, c.description, COALESCE(f.object_key, '')
		FROM celestial_objects c
		LEFT JOIN files f ON c.image_file_id = f.file_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objs []entities.CelestialObject
	for rows.Next() {
		var o entities.CelestialObject
		rows.Scan(&o.ID, &o.Name, &o.Type, &o.Description, &o.ImageURL)
		objs = append(objs, o)
	}
	return objs, nil
}
