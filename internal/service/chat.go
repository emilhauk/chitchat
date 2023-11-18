package service

import (
	"github.com/emilhauk/chitchat/internal/manager"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/pkg/errors"
)

type Chat struct {
	userManager    manager.User
	channelManager manager.Channel
	messageManager manager.Message
}

func NewChatService(userManager manager.User, channelManager manager.Channel, messageManager manager.Message) Chat {
	return Chat{
		userManager:    userManager,
		channelManager: channelManager,
		messageManager: messageManager,
	}
}

func (s Chat) Get(channelUUID string, user model.User) (model.Channel, error) {
	var channel model.Channel
	channel, err := s.channelManager.GetChannelForUser(channelUUID, user.UUID)
	if err != nil {
		return channel, errors.Wrapf(err, "failed to load channel=%s", channelUUID)
	}
	messages, err := s.messageManager.FindMessagesForChannel(channelUUID)
	if err != nil {
		return channel, errors.Wrapf(err, "failed to load messages for channel=%s", channelUUID)
	}
	channel.Messages = messages
	err = s.enhanceMessages(messages, user)
	if err != nil {
		return channel, errors.Wrapf(err, "failed to enhance messages for channel=%s", channelUUID)
	}
	return channel, nil
}

func (s Chat) enhanceMessages(messages []model.Message, user model.User) error {
	userUUIDs := make([]string, 0)
	for i := range messages {
		userUUIDs = append(userUUIDs, messages[i].Sender.UUID)
	}
	users, err := s.userManager.FindAllByUUIDs(userUUIDs...)
	if err != nil {
		return err
	}
	for i := range messages {
		messages[i].Sender = users[messages[i].Sender.UUID]
		messages[i].Direction = model.DirectionIn
		if messages[i].Sender.UUID == user.UUID {
			messages[i].Direction = model.DirectionOut
		}
	}
	return nil
}
