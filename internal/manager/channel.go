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
	FindMember(channelUUID string, userUUID string) (model.Member, error)
}

type Channel struct {
	channelBackend ChannelBackend
}

func NewChannelManager(channelBackend ChannelBackend) Channel {
	return Channel{
		channelBackend: channelBackend,
	}
}

func (m Channel) Create(name string, user model.User) (model.Channel, error) {
	channel := model.Channel{
		UUID:      uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
	}
	err := m.channelBackend.Create(channel)
	if err == nil {
		err = m.channelBackend.AddMember(channel, user, model.RoleAdmin)
	}
	return channel, err
}

func (m Channel) GetChannelListForUser(userUUID string) ([]model.Channel, error) {
	return m.channelBackend.FindAllForUser(userUUID)
}

func (m Channel) GetChannelForUser(channelUUID, userUUID string) (model.Channel, error) {
	return m.channelBackend.FindForUser(channelUUID, userUUID)
}

func (m Channel) GetMemberInfo(channelUUID, userUUID string) (model.Member, error) {
	return m.channelBackend.FindMember(channelUUID, userUUID)
}

func (m Channel) FindByUUID(channelUUID string) (model.Channel, error) {
	return m.channelBackend.FindByUUID(channelUUID)
}

func (m Channel) AddMember(channel model.Channel, user model.User, role model.ChannelRole) error {
	return m.channelBackend.AddMember(channel, user, role)
}
