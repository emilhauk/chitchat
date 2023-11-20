package app

import "errors"

var (
	ErrSessionNotFound              = errors.New("session not found")
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserDeactivated              = errors.New("user is deactivated")
	ErrEmailIsTaken                 = errors.New("email is taken")
	ErrPasswordIncorrect            = errors.New("password is incorrect")
	ErrChannelNotFound              = errors.New("channel not found")
	ErrFieldVerificationNotFound    = errors.New("field verification not found")
	ErrFieldVerificationCodeInvalid = errors.New("field verification code invalid")
	ErrUnsupportedValidationField   = errors.New("unsupported verification field")
	ErrUserHasNoPassword            = errors.New("user has no password")
)
