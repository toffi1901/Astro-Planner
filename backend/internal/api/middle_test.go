package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireRole(t *testing.T) {
	// Создаем тестовый обработчик, который просто отвечает "OK", если до него дошли
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Оборачиваем его в наш Middleware (требуем роль "admin")
	adminOnlyHandler := RequireRole("admin")(nextHandler)

	t.Run("Пользователь с ролью admin", func(t *testing.T) {
		// Создаем фейковый запрос, в контекст которого кладем роль admin
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "role", "admin")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		adminOnlyHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Админ не прошел! Статус: %d", rr.Code)
		}
	})

	t.Run("Пользователь с ролью amateur", func(t *testing.T) {
		// А тут кладем роль amateur
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "role", "amateur")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		adminOnlyHandler.ServeHTTP(rr, req)

		// Сервер должен ответить 403 Forbidden
		if rr.Code != http.StatusForbidden {
			t.Errorf("Обычного юзера пустило! Статус: %d", rr.Code)
		}
	})
}
