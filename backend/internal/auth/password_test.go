package auth

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "super_secret123"
	pepper := "my-pepper"

	// 1. Хеширование
	hash1, err := HashPassword(password, pepper)
	if err != nil {
		t.Fatalf("Ошибка хеширования: %v", err)
	}

	// 2. Проверка, что правильный пароль подходит
	if err := CheckPasswordHash(password, hash1, pepper); err != nil {
		t.Errorf("Правильный пароль не прошел проверку: %v", err)
	}

	// 3. Проверка, что неправильный пароль не подходит
	if err := CheckPasswordHash("wrong_password", hash1, pepper); err == nil {
		t.Error("Неправильный пароль прошел проверку! Это дыра в безопасности.")
	}

	// 4. Проверка соли (два одинаковых пароля должны давать разные хеши)
	hash2, _ := HashPassword(password, pepper)
	if hash1 == hash2 {
		t.Error("Хеши одинаковые! Соль не работает.")
	}
}
