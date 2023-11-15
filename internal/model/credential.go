package model

import "time"

type PasswordCredential struct {
	UserUUID       string
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	LastAssertedAt *time.Time
}
