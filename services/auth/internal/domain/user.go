package domain

import "time"

const (
	RoleUser string = "user"
	RoleAdmin string = "admin"
	RoleCourier string = "courier"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}