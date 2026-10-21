package satellite

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"fmt"
	"time"

	satellite "github.com/joshuaferrara/go-satellite"
)

type PredictUseCase struct {
	Sats  *repository.SatelliteRepository
	Users *repository.PostgresUserRepository
}

func NewPredictUseCase(sats *repository.SatelliteRepository, users *repository.PostgresUserRepository) *PredictUseCase {
	return &PredictUseCase{Sats: sats, Users: users}
}

func (uc *PredictUseCase) Execute(ctx context.Context, userID string, noradID int64, start, end time.Time) ([]entities.Flight, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end_time must be after start_time")
	}
	if end.Sub(start) > 7*24*time.Hour {
		return nil, fmt.Errorf("interval too large (max 7 days)")
	}

	user, err := uc.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Latitude == nil || user.Longitude == nil {
		return nil, fmt.Errorf("user has no coordinates set")
	}
	sat, err := uc.Sats.GetByNoradID(ctx, noradID)
	if err != nil {
		return nil, fmt.Errorf("get satellite: %w", err)
	}
	sgp4sat := satellite.TLEToSat(sat.TLELine1, sat.TLELine2, satellite.GravityWGS84)
	if sgp4sat.Error != 0 {
		return nil, fmt.Errorf("parse TLE: sgp4 error %d", sgp4sat.Error)
	}
	observerLat := *user.Latitude
	observerLon := *user.Longitude

	var passes []entities.Flight
	var currentPass *entities.Flight
	const step = 30 * time.Second

	for t := start; t.Before(end); t = t.Add(step) {
		year, month, day := t.Date()
		hour, min, sec := t.Clock()

		pos, _ := satellite.Propagate(sgp4sat, year, int(month), day, hour, min, sec) //в ECI earth-Centered Inertial
		alt, az := lookAngles(pos, observerLat, observerLon, t)
		if alt > 10.0 {
			// Спутник над горизонтом (> 10 градусов —  минимум для наблюдения)
			if currentPass == nil {
				// Начало нового пролёта
				currentPass = &entities.Flight{
					NoradID:      sat.NoradID,
					SatName:      sat.Name,
					StartTime:    t,
					StartAzimuth: az,
					MaxAltitude:  alt,
				}
			} else {
				if alt > currentPass.MaxAltitude {
					currentPass.MaxAltitude = alt
				}
			}
			currentPass.EndTime = t
			currentPass.EndAzimuth = az
		} else if currentPass != nil {
			// Спутник ушел за горизонт
			passes = append(passes, *currentPass)
			currentPass = nil
		}
	}
	if currentPass != nil {
		passes = append(passes, *currentPass)
	}

	return passes, nil
}

// упрощенная версия
func lookAngles(pos satellite.Vector3, lat, lon float64, t time.Time) (alt, az float64) {
	// Внимание: тут нужна нормальная математика топоцентрических координат.
	_ = pos
	_ = lat
	_ = lon
	_ = t
	return 45.0, 180.0
}
