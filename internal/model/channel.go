package model

import "time"

type ChannelRole = string

const (
	RoleAdmin = "admin"
)

type Channel struct {
	UUID               string
	Name               string
	Messages           []Message
	IsCurrentUserAdmin bool
	CreatedAt          time.Time
	UpdatedAt          *time.Time
}
