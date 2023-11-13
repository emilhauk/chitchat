package model

type Channel struct {
	UUID     string
	Name     string
	Messages []Message
}
