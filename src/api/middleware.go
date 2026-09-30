package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"time"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "duration", time.Since(start))
	})
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateRegister(loginRequest *RegisterInput) error {
	if r.Email == "" || r.Password == "" {
		return errors.New("Login and password must be filled!")
	}

	if !emailPattern.MatchString(r.Email) {
		return errors.New("Incorrect email")
	}
	if len(input.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if !readJSON(w, r, &input) {
		return
	}

	if err := validateRegister(&input); err != nil {
		problem(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.Auth.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		problem(w, http.StatusInternalServerError, "registration failed")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}
