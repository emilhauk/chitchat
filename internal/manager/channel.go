package manager

import (
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

type ChannelBackend interface {
	Create(channel model.Channel) error
	FindByUUID(uuid string) (model.Channel, error)
	FindAllForUser(userUUID string) ([]model.Channel, error)
	FindForUser(channelUUID, userUUID string) (model.Channel, error)
	AddMember(channel model.Channel, user model.User, role model.ChannelRole) error
}

type Channel struct {
	channelBackend ChannelBackend
}

func NewChannelManager(channelBackend ChannelBackend) Channel {
	return Channel{
		channelBackend: channelBackend,
	}
}

func (m Channel) GetChannelListForUser(userUUID string) ([]model.Channel, error) {
	return m.channelBackend.FindAllForUser(userUUID)
}

func (m Channel) GetChannelForUser(channelUUID, userUUID string) (model.Channel, error) {
	return m.channelBackend.FindForUser(channelUUID, userUUID)
}

func (m Channel) isMemberOfChannel(userUUID, channelUUID string) bool {
	for _, c := range channels[userUUID] {
		if c.UUID == channelUUID {
			return true
		}
	}
	return false
}
