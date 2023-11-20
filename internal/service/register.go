package service

import (
	"errors"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"strings"
	"time"
)

var log = config.Logger

type VerificationManager interface {
	CreateAndSendCode(userUUID *string, fieldName string, fieldValue string) (model.FieldVerification, error)
	FindByUUID(uuid string) (model.FieldVerification, error)
}

type UserManager interface {
	Create(user model.User) error
	FindByEmail(email string) (model.User, error)
}

type CredentialManager interface {
	SetPasswordForUser(userUUID, plainPassword string) error
}

type Register struct {
	userManager         UserManager
	verificationManager VerificationManager
	credentialManager   CredentialManager
}

func NewRegisterService(userManager UserManager, verificationManager VerificationManager, credentialManager CredentialManager) Register {
	return Register{
		userManager:         userManager,
		verificationManager: verificationManager,
		credentialManager:   credentialManager,
	}
}

func (s Register) Start(email string) (verification model.FieldVerification, err error) {
	email = strings.ToLower(email)
	_, err = s.userManager.FindByEmail(email)
	if err == nil {
		return verification, app.ErrEmailIsTaken
	} else if !errors.Is(err, app.ErrUserNotFound) {
		return verification, err
	}

	return s.verificationManager.CreateAndSendCode(nil, "email", email)
}

func (s Register) Fulfill(registerRequest model.RegisterRequest) (user model.User, err error) {
	verification, err := s.verificationManager.FindByUUID(registerRequest.VerificationUUID)
	if err != nil {
		return user, err
	}
	var emailVerifiedAt *time.Time
	if verification.Code == registerRequest.Code {
		now := time.Now()
		emailVerifiedAt = &now
	} else if registerRequest.Code == "" && !config.Mail.Enabled {
		log.Warn().Msg("Verification disabled. Allowing unverified email through.")
	} else {
		return user, app.ErrFieldVerificationCodeInvalid
	}

	user = model.User{
		UUID:            uuid.NewString(),
		Name:            registerRequest.Name,
		Email:           verification.FieldValue,
		AvatarUrl:       app.BuildGravatar(verification.FieldValue),
		EmailVerifiedAt: emailVerifiedAt,
		CreatedAt:       time.Now(),
	}
	err = s.userManager.Create(user)
	if err != nil {
		return user, err
	}

	err = s.credentialManager.SetPasswordForUser(user.UUID, registerRequest.PlainPassword)
	if err != nil {
		return user, errors.Join(errors.New("failed to set password for newly created user"), err)
	}
	return user, nil
}
