package repository

import (
	"astro-planner/backend/internal/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SatelliteRepository struct {
	db *pgxpool.Pool
}

type IsatelliteRepository interface {
	Save(ctx context.Context, sat *entities.Satellite) error
	List(ctx context.Context) ([]*entities.Satellite, error)
}

func NewSatelliteRepository(db *pgxpool.Pool) *SatelliteRepository {
	return &SatelliteRepository{db: db}
}

func (r *SatelliteRepository) Save(ctx context.Context, sat *entities.Satellite) error {
	_, err := r.db.Exec(ctx, `
	INSERT INTO satellites (norad_id, satellite_name, satellite_type, tle_line1, tle_line2, last_updated_tle, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (norad_id) DO UPDATE SET 
			satellite_name = EXCLUDED.satellite_name,
			tle_line1 = EXCLUDED.tle_line1,
			tle_line2 = EXCLUDED.tle_line2,
			last_updated_tle = EXCLUDED.last_updated_tle
	`, sat.NoradID, sat.Name, sat.Type, sat.TLELine1, sat.TLELine2, sat.LastUpdatedTLE, sat.IsActive)
	return err
}

func (r *SatelliteRepository) List(ctx context.Context) ([]*entities.Satellite, error) {
	rows, err := r.db.Query(ctx, `
	SELECT norad_id, satellite_name, satellite_type, tle_line1, tle_line2, last_updated_tle, is_active
		FROM satellites
		ORDER BY norad_id `)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var sats []*entities.Satellite
	for rows.Next() {
		sat := &entities.Satellite{}
		err := rows.Scan(&sat.NoradID, &sat.Name, &sat.Type, &sat.TLELine1, &sat.TLELine2, &sat.LastUpdatedTLE, &sat.IsActive)
		if err != nil {
			return nil, err
		}
		sats = append(sats, sat)
	}
	return sats, err
}
