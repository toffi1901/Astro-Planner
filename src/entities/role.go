package entities

type Role int

const (
	RoleGuest     Role = "guest"
	RoleAmateur   Role = "amateur"
	RoleStudent   Role = "student"
	RoleScientist Role = "scientist"
	RoleAdmin     Role = "admin"
)
