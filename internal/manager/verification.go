package manager

import (
	"errors"
	"github.com/emilhauk/chitchat/config"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/google/uuid"
	"time"
)

type VerificationBackend interface {
	Create(verification model.FieldVerification) error
	FindByUUID(uuid string) (model.FieldVerification, error)
}

type Verification struct {
	verificationBackend VerificationBackend
}

func NewVerificationManager(verificationBackend VerificationBackend) Verification {
	return Verification{
		verificationBackend: verificationBackend,
	}
}

func (m Verification) CreateAndSendCode(userUUID *string, fieldName string, fieldValue string) (model.FieldVerification, error) {
	verification := model.FieldVerification{
		UUID:       uuid.NewString(),
		UserUUID:   userUUID,
		FieldName:  fieldName,
		FieldValue: fieldValue,
		CreatedAt:  time.Now(),
	}
	code, err := app.GenerateOTP(6)
	if err != nil {
		return verification, errors.Join(errors.New("failed to generate code"), err)
	}
	verification.Code = code

	err = m.verificationBackend.Create(verification)
	if err != nil {
		return verification, errors.Join(errors.New("failed to create verification"), err)
	}

	if config.Mail.Enabled {
		log.Error().Msg("Sending email is not implemented yet :)")
	} else {
		log.Warn().Msg("Email sending disabled. Allowing users to use whatever they like. Should only be used for testing.")
	}

	return verification, err
}

func (m Verification) FindByUUID(uuid string) (model.FieldVerification, error) {
	return m.verificationBackend.FindByUUID(uuid)
}
