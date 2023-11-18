package manager

import (
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

type MessageBackend interface {
	Create(channelUUID string, message model.Message) error
	FindForChannel(channelUUID string, limit, offset int32) ([]model.Message, error)
}

type Message struct {
	messageBackend MessageBackend
}

func NewMessageManager(messageBackend MessageBackend) Message {
	return Message{
		messageBackend: messageBackend,
	}
}

func (m Message) Send(channel model.Channel, message model.Message) (model.Message, error) {
	message.UUID = uuid.NewString()
	message.Version = 1
	message.SentAt = time.Now()
	message.Direction = model.DirectionOut

	err := m.messageBackend.Create(channel.UUID, message)
	return message, err
}

func (m Message) FindMessagesForChannel(channelUUID string) ([]model.Message, error) {
	return m.messageBackend.FindForChannel(channelUUID, 100, 0)
}
