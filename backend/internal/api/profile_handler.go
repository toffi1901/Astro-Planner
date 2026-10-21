package api

import (
	"astro-planner/backend/internal/usecases/profile"
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
