package app

import "errors"

var (
	ErrSessionNotFound   = errors.New("session not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDeactivated   = errors.New("user is deactivated")
	ErrPasswordIncorrect = errors.New("password is incorrect")
	ErrChannelNotFound   = errors.New("channel not found")
)
