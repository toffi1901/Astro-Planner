package api

import (
	"astro-planner/backend/internal/repository"
	"net/http"
)

type SatelliteHandler struct {
	Repo *repository.SatelliteRepository
}

func (h *SatelliteHandler) List(w http.ResponseWriter, r *http.Request) {
	sats, err := h.Repo.List(r.Context())

	if err != nil {
		problem(w, http.StatusInternalServerError, "failed to list satellites")
		return
	}
	writeJSON(w, http.StatusOK, sats)
}
