package satellite

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	satellite "github.com/joshuaferrara/go-satellite"
)

type PredictUseCase struct {
	Sats  *repository.SatelliteRepository
	Users *repository.PostgresUserRepository
}

var (
	ErrInvalidInterval  = errors.New("end_time must be after start_time")
	ErrIntervalTooLarge = errors.New("interval too large (max 7 days)")
	ErrNoCoordinates    = errors.New("user has no coordinates set")
)

func NewPredictUseCase(sats *repository.SatelliteRepository, users *repository.PostgresUserRepository) *PredictUseCase {
	return &PredictUseCase{Sats: sats, Users: users}
}

func (uc *PredictUseCase) Execute(ctx context.Context, userID string, noradID int64, start, end time.Time) ([]entities.Flight, error) {
	if end.Before(start) {
		return nil, ErrInvalidInterval
	}
	if end.Sub(start) > 7*24*time.Hour {
		return nil, ErrIntervalTooLarge
	}
	user, err := uc.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Latitude == nil || user.Longitude == nil {
		return nil, ErrNoCoordinates
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

func lookAngles(eciPos satellite.Vector3, observerLat, observerLon float64, t time.Time,
) (elevation, azimuth float64) {
	// 1. Считаем юлианскую дату (астрономическое время)
	jday := satellite.JDay(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	// 2. Координаты наблюдателя в градусах
	obsCoords := satellite.LatLong{
		Latitude:  observerLat * math.Pi / 180.0,
		Longitude: observerLon * math.Pi / 180.0,
	}

	// 3. Высота наблюдателя над уровнем моря (в километрах).
	// пока 0 для москвичей
	const observerAlt = 0.0

	// 4. магия: пересчёт ECI в углы обзора
	looks := satellite.ECIToLookAngles(eciPos, obsCoords, observerAlt, jday)

	// Возвращаем высоту над горизонтом и азимут
	return looks.El * 180.0 / math.Pi, looks.Az * 180.0 / math.Pi
}
