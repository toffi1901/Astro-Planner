package satellite

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
)

type GetUseCase struct {
	Sats  *repository.SatelliteRepository
	Users *repository.PostgresUserRepository
}

func NewGetterUseCase(sats *repository.SatelliteRepository, users *repository.PostgresUserRepository) *GetUseCase {
	return &GetUseCase{Sats: sats, Users: users}
}

func (uc *GetUseCase) GetSatellites(ctx context.Context, userID, filter, value string) ([]entities.Satellite, error) {
	user, err := uc.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.IsBlocked {
		return nil, errors.New("user is blocked: no input")
	}
	switch filter {
	case "":
	case "type":
		switch value {
		case "iss", "starlink", "scientific", "weather", "communication", "other":
		default:
			return nil, fmt.Errorf("недопустимый тип спутника: %s", value)
		}

	case "name":
		if strings.TrimSpace(value) == "" {
			return nil, errors.New("поисковый запрос не может быть пустым")
		}

	case "favorites":
		value = userID

	default:
		return nil, fmt.Errorf("неизвестный фильтр: %s", filter)
	}

	sats, err := uc.Sats.GetFaveSatellites(ctx, filter, value)
	if err != nil {
		return nil, fmt.Errorf("get satellites: %w", err)
	}

	return sats, nil
}
