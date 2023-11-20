package database

import (
	"database/sql"
	"errors"
	app "github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type Users struct {
	db *sql.DB

	create            *sql.Stmt
	findByUUID        *sql.Stmt
	findByEmail       *sql.Stmt
	findAllByUUIDs    *sql.Stmt
	findAllByUUIDsSQL string
	setEmail          *sql.Stmt
	setDeactivation   *sql.Stmt
	remove            *sql.Stmt
}

func NewUserStore(db *sql.DB) Users {
	create, err := db.Prepare("INSERT INTO users (uuid, name, email, email_verified_at, created_at) VALUE (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.create")
	}
	findByUUID, err := db.Prepare("SELECT	uuid, name, email, email_verified_at, created_at, last_login_at, deactivated_at, updated_at FROM users WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.findByUUID")
	}
	findByEmail, err := db.Prepare("SELECT uuid, name, email, email_verified_at, created_at, last_login_at, deactivated_at, updated_at FROM users WHERE email = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.findByEmail")
	}
	findAllByUUIDsSQL := "SELECT uuid, name, email, email_verified_at, created_at, last_login_at, deactivated_at, updated_at FROM users WHERE uuid IN (?)"
	findAllByUUIDs, err := db.Prepare(findAllByUUIDsSQL)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.findAllByUUIDs")
	}
	setEmail, err := db.Prepare("UPDATE users SET email = ?, email_verified_at = ?, updated_at = NOW() WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.setEmail")
	}
	setDeactivation, err := db.Prepare("UPDATE users SET deactivated_at = ?, updated_at = NOW() WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.setDeactivation")
	}
	remove, err := db.Prepare("DELETE FROM users WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for users.remove")
	}
	return Users{
		db:                db,
		create:            create,
		findByUUID:        findByUUID,
		findByEmail:       findByEmail,
		findAllByUUIDs:    findAllByUUIDs,
		findAllByUUIDsSQL: findAllByUUIDsSQL,
		setEmail:          setEmail,
		setDeactivation:   setDeactivation,
		remove:            remove,
	}
}

func (s Users) Create(m model.User) error {
	_, err := s.create.Exec(m.UUID, m.Name, m.Email, m.EmailVerifiedAt, m.CreatedAt)
	return err
}

func (s Users) FindByUUID(uuid string) (model.User, error) {
	user, err := s.mapToUser(s.findByUUID.QueryRow(uuid))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user, errors.Join(app.ErrUserNotFound, err)
	}
	return user, err
}

func (s Users) FindByEmail(email string) (model.User, error) {
	user, err := s.mapToUser(s.findByEmail.QueryRow(email))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user, errors.Join(app.ErrUserNotFound, err)
	}
	return user, err
}

func (s Users) FindAllByUUIDs(userUUIDs ...string) (users map[string]model.User, err error) {
	users = map[string]model.User{}
	if len(userUUIDs) == 0 {
		return users, err
	}
	query, args, err := sqlx.In(s.findAllByUUIDsSQL, userUUIDs)
	if err != nil {
		return users, err
	}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return users, err
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return users, err
	}
	for rows.Next() {
		user, err := s.mapToUser(rows)
		if err != nil {
			return users, err
		}
		users[user.UUID] = user
	}
	return users, err
}

func (s Users) SetEmail(uuid, email string, emailVerifiedAt time.Time) error {
	_, err := s.setEmail.Exec(email, emailVerifiedAt, uuid)
	return err
}

func (s Users) SetDeactivation(uuid string, deactivatedAt *time.Time) error {
	_, err := s.setDeactivation.Exec(deactivatedAt, uuid)
	return err
}

func (s Users) Delete(uuid string) error {
	_, err := s.remove.Exec(uuid)
	return err
}

func (s Users) mapToUser(row interface{ Scan(...any) error }) (model.User, error) {
	var (
		uuid            string
		name            string
		email           string
		emailVerifiedAt sql.NullTime
		createdAt       time.Time
		lastLoginAt     sql.NullTime
		deactivatedAt   sql.NullTime
		updatedAt       sql.NullTime
	)
	err := row.Scan(&uuid, &name, &email, &emailVerifiedAt, &createdAt, &lastLoginAt, &deactivatedAt, &updatedAt)
	user := model.User{
		UUID:      uuid,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
		AvatarUrl: app.BuildGravatar(email),
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if deactivatedAt.Valid {
		user.DeactivatedAt = &deactivatedAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = &updatedAt.Time
	}
	return user, err
}
