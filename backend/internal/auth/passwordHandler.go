package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, pepper string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}
func CheckPasswordHash(password, hash, pepper string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+pepper))
}
