package database

import (
	"database/sql"
	"github.com/emilhauk/chitchat/internal/model"
	"time"
)

type Verifications struct {
	db *sql.DB

	create           *sql.Stmt
	findByUUID       *sql.Stmt
	findAllOlderThan *sql.Stmt
	deleteByUUID     *sql.Stmt
}

func NewVerificationsStore(db *sql.DB) Verifications {
	create, err := db.Prepare("INSERT INTO field_verifications (uuid, code, user_uuid, field_name, field_value, created_at) VALUE (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for field_verifications.create")
	}
	findByUUID, err := db.Prepare("SELECT uuid, code, user_uuid, field_name, field_value, created_at FROM field_verifications WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for field_verifications.findByCode")
	}
	findAllOlderThan, err := db.Prepare("SELECT uuid, code, user_uuid, field_name, field_value, created_at FROM field_verifications WHERE created_at < ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for field_verifications.findAllOlderThan")
	}
	deleteByUUID, err := db.Prepare("DELETE FROM field_verifications WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for field_verifications.deleteByUUID")
	}
	return Verifications{
		db:               db,
		create:           create,
		findByUUID:       findByUUID,
		findAllOlderThan: findAllOlderThan,
		deleteByUUID:     deleteByUUID,
	}
}

func (s Verifications) Create(m model.FieldVerification) error {
	_, err := s.create.Exec(m.UUID, m.Code, m.UserUUID, m.FieldName, m.FieldValue, m.CreatedAt)
	return err
}

func (s Verifications) FindByUUID(uuid string) (model.FieldVerification, error) {
	return s.mapToFieldVerification(s.findByUUID.QueryRow(uuid))
}

func (s Verifications) FindAllOlderThan(threshold time.Time) (verifications []model.FieldVerification, err error) {
	verifications = make([]model.FieldVerification, 0)
	rows, err := s.findAllOlderThan.Query(threshold)
	if err != nil {
		return verifications, err
	}
	for rows.Next() {
		verification, err := s.mapToFieldVerification(rows)
		if err != nil {
			return verifications, err
		}
		verifications = append(verifications, verification)
	}
	return verifications, nil
}

func (s Verifications) DeleteByUUID(uuid string) error {
	_, err := s.deleteByUUID.Exec(uuid)
	return err
}

func (s Verifications) mapToFieldVerification(row interface{ Scan(...any) error }) (model.FieldVerification, error) {
	var (
		uuid       string
		code       string
		userUUID   sql.NullString
		fieldName  string
		fieldValue string
		createdAt  time.Time
	)
	err := row.Scan(&uuid, &code, &userUUID, &fieldName, &fieldValue, &createdAt)
	verification := model.FieldVerification{
		UUID:       uuid,
		Code:       code,
		FieldName:  fieldName,
		FieldValue: fieldValue,
		CreatedAt:  createdAt,
	}
	if userUUID.Valid && userUUID.String != "" {
		verification.UserUUID = &userUUID.String
	}
	return verification, err
}
