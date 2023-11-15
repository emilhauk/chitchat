package manager

import (
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

type MessageBackend interface {
	Create(channelUUID string, message model.Message) error
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

	err := m.messageBackend.Create(channel.UUID, message)
	return message, err
}
