package api

import (
	"astro-planner/backend/internal/usecases/profile"
	profileuc "astro-planner/backend/internal/usecases/profile"
	"errors"
	"log"
	"net/http"
)

type ProfileHandler struct {
	ProfileUC *profile.ProfileUseCase
}

// GET /api/profile
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("user_id").(string)
	user, err := h.ProfileUC.GetProfile(r.Context(), userId)
	if err != nil {
		log.Println("ProfileHandler.Get error:", err)
		problem(w, http.StatusInternalServerError, "failed to get profile")
		return
	}
	output := map[string]any{
		"id":        user.ID,
		"email":     user.Email,
		"role":      user.Role,
		"latitude":  user.Latitude,
		"longitude": user.Longitude,
	}
	writeJSON(w, http.StatusOK, output)

}

type UpdateLocationInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// POST /api/profile

func (h *ProfileHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("user_id").(string)
	var input UpdateLocationInput
	if !readJSON(w, r, &input) {
		return
	}
	err := h.ProfileUC.UpdateLocation(r.Context(), userId, input.Latitude, input.Longitude)
	if err != nil {
		if errors.Is(err, profile.ErrInvalidCoordinates) {
			problem(w, http.StatusBadRequest, "coordinates out of range")
			return
		}
		problem(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "update"})
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok || userId == "" {
		problem(w, http.StatusUnauthorized, "no user in context")
		return
	}
	var input ChangePasswordInput
	if !readJSON(w, r, &input) {
		return
	}
	if len(input.NewPassword) < 8 {
		problem(w, http.StatusBadRequest, "new password must be at least 8 characters")
		return
	}
	err := h.ProfileUC.ChangePassword(r.Context(), userId, input.OldPassword, input.NewPassword)
	if err != nil {
		if errors.Is(err, profileuc.ErrWrongOldPassword) {
			problem(w, http.StatusUnauthorized, "incorrect old password")
			return
		}
		if errors.Is(err, profileuc.ErrSamePassword) {
			problem(w, http.StatusBadRequest, "new password must be different")
			return
		}
		problem(w, http.StatusInternalServerError, "failed to change password")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "password changed. please login again"})
}
