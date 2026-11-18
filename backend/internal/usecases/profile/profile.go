package profile

import (
	"astro-planner/backend/internal/auth"
	authcrypto "astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/entities"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type IProfileUseCase interface {
	GetProfile(ctx context.Context, userId string) (*entities.User, error)
}

var ErrInvalidCoordinates = errors.New("invalid coordinates")

type ProfileUseCase struct {
	Users    *repository.PostgresUserRepository
	Sessions *repository.SessionRepository
	Pepper   string
	Mailer   *authcrypto.Mailer
}

func NewProfileUseCase(users *repository.PostgresUserRepository, sessions *repository.SessionRepository, pepper string, mailer *auth.Mailer) *ProfileUseCase {
	return &ProfileUseCase{Users: users,
		Sessions: sessions,
		Pepper:   pepper,
		Mailer:   mailer}
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
	token := uuid.New().String()

	// 2. Сохраняем userId и новый хеш в Redis на 15 минут
	pendingData := repository.PendingPasswordChange{
		UserID:       userId,
		PasswordHash: newHash,
	}
	if err := uc.Sessions.SavePasswordChangeToken(ctx, token, pendingData); err != nil {
		return fmt.Errorf("save password token: %w", err)
	}

	// 3. Отправляем реальное письмо в Mailpit вместо заглушки
	if err := uc.Mailer.SendPasswordConfirmEmail(user.Email, token); err != nil {
		return fmt.Errorf("send confirm email: %w", err)
	}

	return nil
}

// ConfirmPasswordChange вызывается при клике по ссылке из письма
func (uc *ProfileUseCase) ConfirmPasswordChange(ctx context.Context, token string) error {
	// 1. Достаем данные из Redis и сразу удаляем токен
	data, err := uc.Sessions.GetAndDeletePasswordChangeToken(ctx, token)
	if err != nil {
		return fmt.Errorf("get password token: %w", err)
	}

	// 2. Обновляем пароль в PostgreSQL
	if err := uc.Users.UpdatePassword(ctx, data.UserID, data.PasswordHash); err != nil {
		return fmt.Errorf("update password in db: %w", err)
	}

	// 3. Сбрасываем все старые сессии пользователя в Redis
	if err := uc.Sessions.DeleteAllUserSessions(ctx, data.UserID); err != nil {
		log.Println("Ошибка удаления сессий:", err)
	}

	return nil
}

func (uc *ProfileUseCase) GetSessions(ctx context.Context, userId string) ([]*repository.Session, error) {
	return uc.Sessions.ListByUser(ctx, userId)
}

func (uc *ProfileUseCase) RemoveSession(ctx context.Context, jti string) error {
	return uc.Sessions.Delete(ctx, jti)
}
