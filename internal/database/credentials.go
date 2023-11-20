package database

import (
	"database/sql"
	"errors"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"time"
)

type Credentials struct {
	db *sql.DB

	setPassword            *sql.Stmt
	findPasswordByUserUUID *sql.Stmt
}

func NewCredentialStore(db *sql.DB) Credentials {
	setPassword, err := db.Prepare("INSERT INTO password_credentials (user_uuid, password_hash, created_at) VALUE (?, ?, ?) ON DUPLICATE KEY UPDATE password_hash = ?, updated_at = ?")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare statement for password_credentials.setPassword")
	}
	findPasswordByUserUUID, err := db.Prepare("SELECT user_uuid, password_hash, created_at, updated_at, last_asserted_at FROM password_credentials WHERE user_uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare statement for password_credentials.findPasswordByUserUUID")
	}

	return Credentials{
		db:                     db,
		setPassword:            setPassword,
		findPasswordByUserUUID: findPasswordByUserUUID,
	}
}

func (s Credentials) SetPassword(userUUID, hashedPassword string) error {
	_, err := s.setPassword.Exec(userUUID, hashedPassword, time.Now(), hashedPassword, time.Now())
	return err
}

func (s Credentials) FindPasswordByUserUUID(userUUID string) (model.PasswordCredential, error) {
	credential, err := s.mapToPasswordCredential(s.findPasswordByUserUUID.QueryRow(userUUID))
	if errors.Is(err, sql.ErrNoRows) {
		return credential, app.ErrUserHasNoPassword
	}
	return credential, err
}

func (s Credentials) mapToPasswordCredential(row interface{ Scan(...any) error }) (model.PasswordCredential, error) {
	var (
		userUUID       string
		passwordHash   string
		createdAt      time.Time
		updatedAt      sql.NullTime
		lastAssertedAt sql.NullTime
	)
	err := row.Scan(&userUUID, &passwordHash, &createdAt, &updatedAt, &lastAssertedAt)
	credential := model.PasswordCredential{
		UserUUID:     userUUID,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
	if updatedAt.Valid {
		credential.UpdatedAt = &updatedAt.Time
	}
	if lastAssertedAt.Valid {
		credential.LastAssertedAt = &lastAssertedAt.Time
	}
	return credential, err
}
