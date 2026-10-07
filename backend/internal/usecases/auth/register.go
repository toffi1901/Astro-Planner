package auth

import (
	"astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
)

type RegisterUseCase struct {
	Pepper string
	Users  *repository.PostgresUserRepository
}

func NewRegisterUseCase(users *repository.PostgresUserRepository, pepper string) *RegisterUseCase {
	return &RegisterUseCase{Users: users, Pepper: pepper}
}

func (usecase *RegisterUseCase) Execute(ctx context.Context, email, password string) (*entities.User, error) {
	hash, err := auth.HashPassword(password, usecase.Pepper)
	if err != nil {
		return nil, err
	}
	user := &entities.User{
		Email:        email,
		PasswordHash: hash,
		Role:         "amateur",
	}
	if err := usecase.Users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
