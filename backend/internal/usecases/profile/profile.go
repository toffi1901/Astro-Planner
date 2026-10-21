package profile

import (
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
)

type IProfileUseCase interface {
	GetProfile(ctx context.Context, userId string) (*entities.User, error)
}

var ErrInvalidCoordinates = errors.New("invalid coordinates")

type ProfileUseCase struct {
	Users *repository.PostgresUserRepository
}

func NewProfileUseCase(users *repository.PostgresUserRepository) *ProfileUseCase {
	return &ProfileUseCase{Users: users}
}

func (uc *ProfileUseCase) GetProfile(ctx context.Context, userId string) (*entities.User, error) {
	return uc.Users.GetByID(ctx, userId)
}

func (uc *ProfileUseCase) UpdateLocation(ctx context.Context, userId string, latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return ErrInvalidCoordinates
	}
	return uc.Users.UpdateCoordinates(ctx, userId, latitude, longitude)
}
