package database

import (
	"database/sql"
	"github.com/emilhauk/chitchat/internal/model"
	"time"
)

type Credentials struct {
	db *sql.DB

	createPassword         *sql.Stmt
	findPasswordByUserUUID *sql.Stmt
}

func NewCredentialStore(db *sql.DB) Credentials {
	createPassword, err := db.Prepare("INSERT INTO password_credentials (user_uuid, password_hash, created_at) VALUE (?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare statement for password_credentials.createPassword")
	}
	findPasswordByUserUUID, err := db.Prepare("SELECT user_uuid, password_hash, created_at, updated_at, last_asserted_at FROM password_credentials WHERE user_uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to prepare statement for password_credentials.findPasswordByUserUUID")
	}

	return Credentials{
		db:                     db,
		createPassword:         createPassword,
		findPasswordByUserUUID: findPasswordByUserUUID,
	}
}

func (s Credentials) FindPasswordByUserUUID(userUUID string) (model.PasswordCredential, error) {
	return s.mapToPasswordCredential(s.findPasswordByUserUUID.QueryRow(userUUID))
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
