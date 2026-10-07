package api

import (
	authuc "astro-planner/backend/internal/usecases/auth"
	"errors"
	"net/http"
	"regexp"
)

type AuthHandler struct {
	RegisterUC *authuc.RegisterUseCase
	LoginUC    *authuc.LoginUseCase
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutput struct {
	ID int `json:"id"`
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
