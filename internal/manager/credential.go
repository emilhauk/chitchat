package manager

import (
	"errors"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type CredentialBackend interface {
	SetPassword(userUUID, hashedPassword string) error
	FindPasswordByUserUUID(userUUID string) (model.PasswordCredential, error)
}

type Credential struct {
	credentialBackend CredentialBackend
}

func NewCredentialManager(verificationBackend CredentialBackend) Credential {
	return Credential{
		credentialBackend: verificationBackend,
	}
}

func (m Credential) SetPasswordForUser(userUUID, plainPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Join(errors.New("failed to hash password"), err)
	}
	err = m.credentialBackend.SetPassword(userUUID, string(passwordHash))
	if err != nil {
		return errors.Join(errors.New("failed to store hashed password"), err)
	}
	return nil
}

func (m Credential) CheckPasswordForUser(userUUID, plainPassword string) (isValid bool, err error) {
	credential, err := m.credentialBackend.FindPasswordByUserUUID(userUUID)
	if err != nil {
		if errors.Is(err, app.ErrUserHasNoPassword) {
			return false, nil
		}
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(plainPassword))
	if err != nil {
		log.Debug().Err(err).Msg("BCrypt password comparison failed")
		return false, nil
	}
	return true, nil
}
