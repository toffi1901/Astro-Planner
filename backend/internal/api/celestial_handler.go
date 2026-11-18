package api

import (
	"astro-planner/backend/internal/repository"
	"astro-planner/backend/internal/storage"
	"net/http"
)

type CelestialHandler struct {
	Repo *repository.CelestialRepo
	S3   *storage.S3Storage
}

func (h *CelestialHandler) List(w http.ResponseWriter, r *http.Request) {
	objs, err := h.Repo.List(r.Context())
	if err != nil {
		problem(w, http.StatusInternalServerError, "failed to get objects")
		return
	}

	for i := range objs {
		if objs[i].ImageURL != "" {
			objs[i].ImageURL = h.S3.GetPublicURL(objs[i].ImageURL)
		}
	}
	writeJSON(w, http.StatusOK, objs)
}
