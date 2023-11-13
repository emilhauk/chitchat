package manager

import (
	"errors"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

var channels = map[string][]model.Channel{
	"0-1-3": {
		{
			UUID:     "1-13-37",
			Name:     "Project BEMRØ",
			Messages: make([]model.Message, 0),
		}, {
			UUID:     "1-13-38",
			Name:     "Haukeland",
			Messages: make([]model.Message, 0),
		},
	},
}

var (
	NoSuchChannelError      = errors.New("no such channel")
	NotMemberOfChannelError = errors.New("not member channel")
)

type Channel struct {
}

func (m Channel) GetChannelsForUser(userUUID string) ([]model.Channel, error) {
	return channels[userUUID], nil
}

func (m Channel) SendMessage(channelUUID string, message model.Message) (model.Message, error) {
	if !m.isMemberOfChannel(message.Sender.UUID, channelUUID) {
		return message, NoSuchChannelError
	}
	message.UUID = uuid.NewString()
	message.SentAt = time.Now()

	for i, c := range channels[message.Sender.UUID] {
		if c.UUID != channelUUID {
			continue
		}
		channels[message.Sender.UUID][i].Messages = append(c.Messages, message)
	}

	return message, nil
}

func (m Channel) isMemberOfChannel(userUUID, channelUUID string) bool {
	for _, c := range channels[userUUID] {
		if c.UUID == channelUUID {
			return true
		}
	}
	return false
}

func (m Channel) GetChannelForUser(channelUUID, userUUID string) (model.Channel, error) {
	for _, c := range channels[userUUID] {
		if c.UUID == channelUUID {
			return c, nil
		}
	}
	return model.Channel{}, NoSuchChannelError
}
