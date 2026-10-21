package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateToken(userID, role, jti, secret string) (*TokenPair, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"jti":     jti,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	resfreshClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"jti":     jti,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, resfreshClaims).SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: token, RefreshToken: refreshToken}, nil
}

func ValidateToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
