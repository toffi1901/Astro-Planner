package auth

import (
	authcrypto "astro-planner/backend/internal/auth"
	"astro-planner/backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshUseCase struct {
	Sessions  *repository.SessionRepository
	JWTSecret string
}

func NewRefreshUseCase(session *repository.SessionRepository, jwtSecret string) *RefreshUseCase {
	return &RefreshUseCase{Sessions: session, JWTSecret: jwtSecret}
}

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (*authcrypto.TokenPair, error) {
	claims, err := authcrypto.ValidateToken(refreshToken, uc.JWTSecret)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if claims["type"] != "refresh" {
		return nil, ErrInvalidRefreshToken
	}

	oldJti, _ := claims["jti"].(string)
	session, err := uc.Sessions.Get(ctx, oldJti)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	_ = uc.Sessions.Delete(ctx, oldJti)
	newJti := uuid.New().String()
	newSession := repository.Session{
		UserID:    session.UserID,
		Role:      session.Role,
		IPadress:  session.IPadress,
		UserAgent: session.UserAgent,
		CreatedAt: time.Now(),
	}
	if err := uc.Sessions.Create(ctx, newJti, newSession, 7*24*time.Hour); err != nil {
		return nil, err
	}
	return authcrypto.GenerateToken(session.UserID, session.Role, newJti, uc.JWTSecret)
}
