package entities

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Role         string
	IsBlocked    bool
	RegisterDate time.Time
}
