package model

import "time"

type User struct {
	UUID            string
	Name            string
	Email           string
	AvatarUrl       string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	LastLoginAt     *time.Time
	DeactivatedAt   *time.Time
	UpdatedAt       *time.Time
}

type FieldVerification struct {
	UUID       string
	Code       string
	UserUUID   *string
	FieldName  string
	FieldValue string
	CreatedAt  time.Time
}

type RegisterRequest struct {
	VerificationUUID string
	Code             string
	Name             string
	PlainPassword    string
}
