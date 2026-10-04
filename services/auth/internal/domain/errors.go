package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("User not found.")
	ErrUserAlreadyExist    = errors.New("User already exist.")
	ErrInvalidEmail        = errors.New("Invalid email.")
	ErrInvalidCredentials  = errors.New("Invalid сredentials.")
	ErrInvalidToken        = errors.New("Invalid token.")
	ErrInvalidAdminCode    = errors.New("Invalid admin code.")
)