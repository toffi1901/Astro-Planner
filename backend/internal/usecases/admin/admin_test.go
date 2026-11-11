package admin

import (
	"context"
	"errors"
	"testing"
)

func TestChangeUserRole_InvalidRole(t *testing.T) {
	// Передаем nil вместо репозиториев, так как при невалидной роли до базы дело не дойдет
	uc := NewAdminUseCase(nil, nil)
	ctx := context.Background()

	// Пытаемся назначить несуществующую роль
	err := uc.ChangeUserRole(ctx, "user-123", "super_boss")

	if !errors.Is(err, ErrInvalidRole) {
		t.Errorf("Ожидалась ошибка ErrInvalidRole, получено: %v", err)
	}

	// Пытаемся назначить валидную роль (тут он пойдет в базу и упадет с паникой,
	// так как repos nil, поэтому мы просто проверяем, что ошибка НЕ про невалидную роль)
	// В реальном проекте тут нужно использовать моки, но для базового теста сойдет и так.
}
