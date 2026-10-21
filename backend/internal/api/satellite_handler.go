package api

import (
	"astro-planner/backend/internal/repository"
	"errors"
	"net/http"
	"strconv"
	"time"

	satuc "astro-planner/backend/internal/usecases/satellite"

	"github.com/go-chi/chi/v5"
)

type SatelliteHandler struct {
	Repo      *repository.SatelliteRepository
	PredictUC *satuc.PredictUseCase
}

type PredictRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func (h *SatelliteHandler) List(w http.ResponseWriter, r *http.Request) {
	sats, err := h.Repo.List(r.Context())

	if err != nil {
		problem(w, http.StatusInternalServerError, "failed to list satellites")
		return
	}
	writeJSON(w, http.StatusOK, sats)
}

func (h *SatelliteHandler) Predict(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)
	if userID == "" {
		problem(w, http.StatusUnauthorized, "no user")
		return
	}

	noradStr := chi.URLParam(r, "norad_id")
	noradID, err := strconv.ParseInt(noradStr, 10, 64)
	if err != nil {
		if errors.Is(err, repository.ErrSatelliteNotFound) {
			problem(w, http.StatusBadRequest, "satellite not found")
			return
		}
	}

	var input PredictRequest
	if !readJSON(w, r, &input) {
		return
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		problem(w, http.StatusBadRequest, "start_time and end_time are required")
		return
	}
	passes, err := h.PredictUC.Execute(r.Context(), userID, noradID, input.StartTime, input.EndTime)
	if err != nil {
		if errors.Is(err, repository.ErrSatelliteNotFound) {
			problem(w, http.StatusNotFound, "satellite not found")
			return
		}
		if errors.Is(err, satuc.ErrInvalidInterval) ||
			errors.Is(err, satuc.ErrIntervalTooLarge) ||
			errors.Is(err, satuc.ErrNoCoordinates) {
			problem(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"satellite_id": noradID,
		"passes":       passes,
		"count":        len(passes),
	})
}
