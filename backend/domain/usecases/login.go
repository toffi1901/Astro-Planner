package usecases

import "net/http"

type LoginInput struct {
	Login    string `json:"login"`
	Password string `json: "pqssword"`
}

func (a *app) login(w http.ResponseWriter)
