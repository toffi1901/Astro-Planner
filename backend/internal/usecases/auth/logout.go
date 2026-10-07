package auth

import (
	"astro-planner/backend/internal/repository"
	"context"
)

type LogoutUseCase struct {
	Sessions *repository.SessionRepository
}

func NewLogoutUseCase(sessions *repository.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{Sessions: sessions}
}

func (usecase *LogoutUseCase) Execute(ctx context.Context, jti string) error {
	return usecase.Sessions.Delete(ctx, jti)
}
