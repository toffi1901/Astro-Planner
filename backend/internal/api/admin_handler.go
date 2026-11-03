package api

import (
	adminuc "astro-planner/backend/internal/usecases/admin"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	AdminUC *adminuc.AdminUseCase
}

type ChangeRoleInput struct {
	Role string `json:"role"`
}

func (h *AdminHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	var input ChangeRoleInput
	if !readJSON(w, r, &input) {
		return
	}
	err := h.AdminUC.ChangeUserRole(r.Context(), userId, input.Role)
	if err != nil {
		log.Println(" ОШИБКА СМЕНЫ РОЛИ:", err)
		if errors.Is(err, adminuc.ErrInvalidRole) {
			problem(w, http.StatusBadRequest, "invalid role")
			return
		}
		problem(w, http.StatusInternalServerError, "failed to change role")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "role updated"})
}

type BlockInput struct {
	IsBlocked bool `json:"is_blocked"`
}

func (h *AdminHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	var input BlockInput
	if !readJSON(w, r, &input) {
		return
	}
	if err := h.AdminUC.SetBlockStatus(r.Context(), userId, input.IsBlocked); err != nil {
		problem(w, http.StatusInternalServerError, "failed to change block status")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "block status updated"})
}
