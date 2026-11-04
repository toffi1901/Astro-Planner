package profile

import (
	authcrypto "astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
	"log"
)

type IProfileUseCase interface {
	GetProfile(ctx context.Context, userId string) (*entities.User, error)
}

var ErrInvalidCoordinates = errors.New("invalid coordinates")

type ProfileUseCase struct {
	Users    *repository.PostgresUserRepository
	Sessions *repository.SessionRepository
	Pepper   string
}

func NewProfileUseCase(users *repository.PostgresUserRepository, sessions *repository.SessionRepository, pepper string) *ProfileUseCase {
	return &ProfileUseCase{Users: users,
		Sessions: sessions,
		Pepper:   pepper}
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

var (
	ErrWrongOldPassword = errors.New("incorrect old password")
	ErrSamePassword     = errors.New("new password must be different")
)

func (uc *ProfileUseCase) ChangePassword(ctx context.Context, userId, oldPass, newPass string) error {
	user, err := uc.Users.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	if err := authcrypto.CheckPasswordHash(oldPass, user.PasswordHash, uc.Pepper); err != nil {
		return ErrWrongOldPassword
	}
	if oldPass == newPass {
		return ErrSamePassword
	}
	newHash, err := authcrypto.HashPassword(newPass, uc.Pepper)
	if err != nil {
		return err
	}
	if err := uc.Users.UpdatePassword(ctx, userId, newHash); err != nil {
		log.Println("Ошибка удаления сессий:", err)
	}
	log.Printf("[EMAIL] Пользователю %s отправлено письмо: 'Ваш пароль был успешно изменён.'", user.Email) //заглушка
	return nil
}

func (uc *ProfileUseCase) GetSessions(ctx context.Context, userId string) ([]*repository.Session, error) {
	return uc.Sessions.ListByUser(ctx, userId)
}

func (uc *ProfileUseCase) RemoveSession(ctx context.Context, jti string) error {
	return uc.Sessions.Delete(ctx, jti)
}
