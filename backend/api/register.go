package api

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	authuc "github.com/toffi_1901/astro-planner/backend/auth"
)

type AuthHandler struct {
	RegisterUC *authuc.RegisterUseCase
	LoginUC    *authuc.LoginUseCase
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userDB struct {
	ID    int
	Email string
	Salt  []byte
	Hash  []byte
	Role  Role
}

type UserOutput struct {
	ID int `json:"id"`
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "duration", time.Since(start))
	})
}

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateRegister(input *RegisterInput) error {
	if input.Email == "" || input.Password == "" {
		return errors.New("Login and password must be filled!")
	}

	if !emailPattern.MatchString(input.Email) {
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

	user, err := h.RegisterUC.Execute(r.Context(), input.Email, input.Password)
	if err != nil {
		problem(w, http.StatusInternalServerError, "registration failed")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}
