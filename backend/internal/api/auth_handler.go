package api

import (
	authuc "astro-planner/backend/internal/usecases/auth"
	"errors"
	"net"
	"net/http"
	"regexp"
)

type AuthHandler struct {
	RegisterUC *authuc.RegisterUseCase
	LoginUC    *authuc.LoginUseCase
	LogoutUC   *authuc.LogoutUseCase
	RefreshUC  *authuc.RefreshUseCase
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutput struct {
	ID int `json:"id"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if !readJSON(w, r, &input) {
		return
	}
	if input.Email == "" || input.Password == "" {
		problem(w, http.StatusBadRequest, "email and password are required")
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	token, err := h.LoginUC.Execute(r.Context(), input.Email, input.Password, net.ParseIP(ip), r.UserAgent())
	if err != nil {
		if errors.Is(err, authuc.ErrTooManyAttempts) {
			problem(w, http.StatusTooManyRequests, "too many attempts")
			return
		}
		problem(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})
	writeJSON(w, http.StatusOK, map[string]string{
		"access_token": token.AccessToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	jti, _ := r.Context().Value("jti").(string)
	if jti == "" {
		problem(w, http.StatusUnauthorized, "no session")
		return
	}
	if err := h.LogoutUC.Execute(r.Context(), jti); err != nil {
		problem(w, http.StatusNoContent, "logout failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
