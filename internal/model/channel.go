package model

import "time"

type Channel struct {
	UUID      string
	Name      string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt *time.Time
}
