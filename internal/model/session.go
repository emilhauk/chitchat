package model

import "time"

type Session struct {
	ID         string
	UserUUID   string
	CreatedAt  time.Time
	LastSeenAt *time.Time
}
