package profile

import (
	"context"
	"errors"
	"testing"
)

func TestUpdateLocationValidation(t *testing.T) {
	// Нам не нужна база данных, потому что валидация происходит ДО запроса в БД
	uc := &ProfileUseCase{}
	ctx := context.Background()

	// Тест 1: Широта больше 90 (должна быть ошибка)
	err := uc.UpdateLocation(ctx, "user-id", 95.0, 37.0)
	if !errors.Is(err, ErrInvalidCoordinates) {
		t.Errorf("Ожидалась ошибка ErrInvalidCoordinates, получено: %v", err)
	}

	// Тест 2: Долгота больше 180 (должна быть ошибка)
	err = uc.UpdateLocation(ctx, "user-id", 55.0, 185.0)
	if !errors.Is(err, ErrInvalidCoordinates) {
		t.Errorf("Ожидалась ошибка ErrInvalidCoordinates, получено: %v", err)
	}

	// (Тест на успешное сохранение пропускаем, так как там нужен реальный UsersRepo)
}
