package model

import "time"

type Direction = string

const (
	DirectionIn  = "in"
	DirectionOut = "out"
)

type Message struct {
	UUID      string
	Sender    User
	Content   string `json:"content"`
	Direction Direction
	Version   uint32
	SentAt    time.Time
	DeletedAt *time.Time
	UpdatedAt *time.Time
}
