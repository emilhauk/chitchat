package model

import "time"

type Message struct {
	UUID    string
	Sender  User
	Content string `json:"content"`
	SentAt  time.Time
}
