package api

import (
	"net/http"
)

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		problem(w, http.StatusUnauthorized, "refresh token missing")
		return
	}
	pair, err := h.RefreshUC.Execute(r.Context(), cookie.Value)
	if err != nil {
		problem(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60,
	})
	writeJSON(w, http.StatusOK, map[string]string{
		"access_token": pair.AccessToken,
	})

}
