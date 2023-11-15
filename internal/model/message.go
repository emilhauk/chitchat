package model

import "time"

type Message struct {
	UUID      string
	Sender    User
	Content   string `json:"content"`
	Version   uint32
	SentAt    time.Time
	DeletedAt *time.Time
	UpdatedAt *time.Time
}
