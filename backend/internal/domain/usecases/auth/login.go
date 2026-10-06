package usecases

import (
	"astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/domain/repository"
	"context"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidData     = errors.New("invalid credentials")
	ErrTooManyAttempts = errors.New("too many attempts")
)

type LoginUseCase struct {
	Users     *repository.PostgresUserRepository
	Sessions  *repository.SessionRepository
	Limiter   *repository.LimitsRepository
	Pepper    string
	JWTSecret string
}

func NewLoginUseCase(users *repository.PostgresUserRepository,
	session *repository.SessionRepository,
	limiter *repository.LimitsRepository,
	pepper, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{Users: users,
		Sessions:  session,
		Limiter:   limiter,
		Pepper:    pepper,
		JWTSecret: jwtSecret}
}

func (usecase *LoginUseCase) Execute(ctx context.Context, email, password string, ip net.IP, EnvironmentInfo string) (string, error) {
	blocked, err := usecase.Limiter.IsBlocked(ctx, ip.String())
	if err != nil {
		return "", err
	}
	if blocked {
		return "", ErrTooManyAttempts
	}
	user, err := usecase.Users.GetByEmail(ctx, email)
	if err != nil {
		usecase.Limiter.Increment(ctx, ip)
		return "", ErrInvalidData
	}

	if err := auth.CheckPasswordHash(password, user.PasswordHash, usecase.Pepper); err != nil {
		usecase.Limiter.Increment(ctx, ip)
		return "", ErrInvalidData
	}
	usecase.Limiter.Reset(ctx, ip)
	jti := uuid.New().String()
	session := repository.Session{
		UserID:          user.ID,
		Role:            user.Role,
		IPadress:        ip.String(),
		EnvironmentInfo: EnvironmentInfo,
		CreatedAt:       time.Now(),
	}

	if err := usecase.Sessions.Create(ctx, jti, session, 24*time.Hour); err != nil {
		return "", err
	}
	return auth.GenerateToken(session.UserID, user.Role, jti, usecase.JWTSecret)
}
