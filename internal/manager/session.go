package manager

import (
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

type SessionBackend interface {
	Create(model model.Session) error
	FindByID(id string) (model.Session, error)
	SetLastSeenAt(id string, lastSeenAt time.Time) error
	Delete(id string) error
}

type Session struct {
	sessionBackend SessionBackend
}

func NewSessionManager(sessionBackend SessionBackend) Session {
	return Session{
		sessionBackend: sessionBackend,
	}
}

func (m Session) CreateSession(userUUID string) (model.Session, error) {
	session := model.Session{
		ID:        uuid.NewString(),
		UserUUID:  userUUID,
		CreatedAt: time.Now(),
	}

	err := m.sessionBackend.Create(session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func (m Session) FindByID(id string) (model.Session, error) {
	return m.sessionBackend.FindByID(id)
}

func (m Session) SetLastSeenAt(id string, lastSeenAt time.Time) error {
	return m.sessionBackend.SetLastSeenAt(id, lastSeenAt)
}

func (m Session) Delete(id string) error {
	return m.sessionBackend.Delete(id)
}
