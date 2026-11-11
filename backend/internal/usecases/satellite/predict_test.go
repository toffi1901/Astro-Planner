package satellite

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAzimuthToDirection(t *testing.T) {
	// Таблица тестовых случаев
	tests := []struct {
		azimuth  float64
		expected string
	}{
		{0, "С"},
		{90, "В"},
		{180, "Ю"},
		{270, "З"},
		{45, "СВ"},
		{225, "ЮЗ"},
	}

	for _, tt := range tests {
		result := azimuthToDirection(tt.azimuth)
		if result != tt.expected {
			t.Errorf("Для азимута %v ожидалось %s, получено %s", tt.azimuth, tt.expected, result)
		}
	}
}

func TestPredictExecute_InvalidIntervals(t *testing.T) {
	// Опять же, repos nil, так как валидация времени идет первой
	uc := NewPredictUseCase(nil, nil)
	ctx := context.Background()

	now := time.Now()

	// Тест 1: Конец раньше начала
	_, err := uc.Execute(ctx, "user-id", 25544, now, now.Add(-1*time.Hour))
	if !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("Ожидалась ошибка ErrInvalidInterval, получено: %v", err)
	}

	// Тест 2: Интервал больше 7 дней
	_, err = uc.Execute(ctx, "user-id", 25544, now, now.Add(8*24*time.Hour))
	if !errors.Is(err, ErrIntervalTooLarge) {
		t.Errorf("Ожидалась ошибка ErrIntervalTooLarge, получено: %v", err)
	}
}
