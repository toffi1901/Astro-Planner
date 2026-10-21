package auth

import "testing"

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "test-secret-key"
	userID := "user-123"
	role := "admin"
	jti := "jti-456"

	// 1. Генерируем пару токенов
	pair, err := GenerateToken(userID, role, jti, secret)
	if err != nil {
		t.Fatalf("Ошибка генерации токена: %v", err)
	}

	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("Токены пустые!")
	}

	// 2. Валидируем Access токен
	claims, err := ValidateToken(pair.AccessToken, secret)
	if err != nil {
		t.Fatalf("Ошибка валидации правильного токена: %v", err)
	}

	// 3. Проверяем содержимое
	if claims["user_id"] != userID {
		t.Errorf("Ожидался user_id=%s, получен %v", userID, claims["user_id"])
	}
	if claims["type"] != "access" {
		t.Errorf("Ожидался type=access, получен %v", claims["type"])
	}

	// 4. Проверяем, что чужой секрет не подходит
	_, err = ValidateToken(pair.AccessToken, "wrong-secret")
	if err == nil {
		t.Error("Токен прошел проверку с неправильным секретом!")
	}
}
