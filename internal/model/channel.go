package model

import (
	"time"
)

type ChannelRole = string

const (
	RoleAdmin ChannelRole = "admin"
)

type Channel struct {
	UUID               string
	Name               string
	Messages           []Message
	IsCurrentUserAdmin bool
	InvitationURL      string
	CreatedAt          time.Time
	UpdatedAt          *time.Time
}

type Member struct {
	ChannelUUID string
	UserUUID    string
	Role        ChannelRole
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
