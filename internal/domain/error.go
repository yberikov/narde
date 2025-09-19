package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidMove        = errors.New("invalid move")
	ErrOpponentNotFound   = errors.New("opponent not found")
)
