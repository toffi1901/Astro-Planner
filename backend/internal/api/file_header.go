package api

import (
	fileuc "astro-planner/backend/internal/usecases/file"
	"net/http"
)

type FileHandler struct {
	FileUC *fileuc.FileUseCase
}

// POST /api/files/upload
func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)
	if userID == "" {
		problem(w, http.StatusUnauthorized, "no user")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		problem(w, http.StatusBadRequest, "file too large")
		return
	}

	// достаем файл (ключ "file" — под этим именем фронт отправляет файл)
	file, header, err := r.FormFile("file")
	if err != nil {
		problem(w, http.StatusBadRequest, "missing file")
		return
	}
	defer file.Close()

	f, url, err := h.FileUC.UploadFile(r.Context(), userID, file, header)
	if err != nil {
		problem(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"file_id": f.ID,
		"url":     url,
	})
}
