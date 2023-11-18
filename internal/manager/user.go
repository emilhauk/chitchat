package manager

import (
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var log = config.Logger

type UserBackend interface {
	Create(model.User) error
	FindByUUID(uuid string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindAllByUUIDs(userUUIDs ...string) (map[string]model.User, error)
}

type CredentialBackend interface {
	FindPasswordByUserUUID(uuid string) (model.PasswordCredential, error)
}

type User struct {
	userBackend       UserBackend
	credentialBackend CredentialBackend
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

func (m User) FindByEmailAndPlainPassword(email, plainPassword string) (model.User, error) {
	user, err := m.userBackend.FindByEmail(email)
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
