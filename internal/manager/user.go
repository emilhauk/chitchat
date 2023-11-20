package manager

import (
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var log = config.Logger

type UserBackend interface {
	Create(model.User) error
	FindByUUID(uuid string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindAllByUUIDs(userUUIDs ...string) (map[string]model.User, error)
	SetEmail(uuid, email string, emailVerifiedAt time.Time) error
}

type User struct {
	userBackend       UserBackend
	credentialBackend CredentialBackend
}

func (m User) Create(user model.User) error {
	return m.userBackend.Create(user)
}

func NewUserManager(userBackend UserBackend, credentialBackend CredentialBackend) User {
	return User{
		userBackend:       userBackend,
		credentialBackend: credentialBackend,
	}
}

func (m User) FindByUUID(uuid string) (model.User, error) {
	return m.userBackend.FindByUUID(uuid)
}

func (m User) FindByEmail(email string) (model.User, error) {
	return m.userBackend.FindByEmail(email)
}

func (m User) FindByEmailAndPlainPassword(email, plainPassword string) (user model.User, err error) {
	user, err = m.userBackend.FindByEmail(email)
	if err != nil {
		return user, err
	}

	password, err := m.credentialBackend.FindPasswordByUserUUID(user.UUID)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password.PasswordHash), []byte(plainPassword))
	if err != nil {
		log.Debug().Err(err).Msg("Password comparison failed")
		return user, app.ErrPasswordIncorrect
	}

	return user, err
}

func (m User) FindAllByUUIDs(userUUIDs ...string) (map[string]model.User, error) {
	return m.userBackend.FindAllByUUIDs(userUUIDs...)
}
