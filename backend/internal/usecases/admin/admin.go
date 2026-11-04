package admin

import (
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
)

var ErrInvalidRole = errors.New("invalid role")

type AdminUseCase struct {
	Users    *repository.PostgresUserRepository
	Sessions *repository.SessionRepository
}

func NewAdminUseCase(users *repository.PostgresUserRepository, sessions *repository.SessionRepository) *AdminUseCase {
	return &AdminUseCase{Users: users, Sessions: sessions}
}

func (uc *AdminUseCase) ChangeUserRole(ctx context.Context, userId, newRole string) error {
	validRoles := map[string]bool{
		"guest": true, "amateur": true, "student": true, "scientist": true, "admin": true,
	}
	if !validRoles[newRole] {
		return ErrInvalidRole
	}
	if err := uc.Users.UpdateRole(ctx, userId, newRole); err != nil {
		return err
	}
	_ = uc.Sessions.DeleteAllUserSessions(ctx, userId)
	return nil
}

func (uc *AdminUseCase) SetBlockStatus(ctx context.Context, userId string, isBlocked bool) error {
	if err := uc.Users.SetBlockUser(ctx, userId, isBlocked); err != nil {
		return err
	}
	if isBlocked {
		_ = uc.Sessions.DeleteAllUserSessions(ctx, userId)
	}
	return nil
}

func (uc *AdminUseCase) RemoveAllUserSessions(ctx context.Context, userId string) error {
	return uc.Sessions.DeleteAllUserSessions(ctx, userId)
}
