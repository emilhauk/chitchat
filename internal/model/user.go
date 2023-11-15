package model

import "time"

type User struct {
	UUID            string
	Name            string
	Email           string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	LastLoginAt     *time.Time
	DeactivatedAt   *time.Time
	UpdatedAt       *time.Time
}
