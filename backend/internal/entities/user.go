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
	Latitude     *float64
	Longitude    *float64
}
