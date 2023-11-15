package database

import (
	"database/sql"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/rs/zerolog/log"
	"time"
)

type Sessions struct {
	db *sql.DB

	create           *sql.Stmt
	findById         *sql.Stmt
	updateLastSeenAt *sql.Stmt
	delete           *sql.Stmt
}

func NewSessionStore(db *sql.DB) Sessions {
	create, err := db.Prepare("INSERT INTO sessions (id, user_uuid, created_at, last_seen_at) VALUE (?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for sessions.create")
	}
	findById, err := db.Prepare("SELECT id, user_uuid, created_at, last_seen_at FROM sessions WHERE id = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for sessions.findById")
	}
	updateLastSeenAt, err := db.Prepare("UPDATE sessions SET last_seen_at = ? WHERE id = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for sessions.updateLastSeenAt")
	}
	remove, err := db.Prepare("DELETE FROM sessions WHERE id = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for sessions.remove")
	}
	return Sessions{
		db:               db,
		create:           create,
		findById:         findById,
		updateLastSeenAt: updateLastSeenAt,
		delete:           remove,
	}
}

func (s Sessions) Create(m model.Session) error {
	_, err := s.create.Exec(m.ID, m.UserUUID, m.CreatedAt, m.LastSeenAt)
	return err
}

func (s Sessions) FindByID(id string) (model.Session, error) {
	return s.mapToSession(s.findById.QueryRow(id))
}

func (s Sessions) SetLastSeenAt(id string, lastSeenAt time.Time) error {
	_, err := s.updateLastSeenAt.Exec(lastSeenAt, id)
	return err
}

func (s Sessions) Delete(id string) error {
	_, err := s.delete.Exec(id)
	return err
}

func (s Sessions) mapToSession(row interface{ Scan(...any) error }) (model.Session, error) {
	var (
		id         string
		userUUID   string
		createdAt  time.Time
		lastSeenAt sql.NullTime
	)

	err := row.Scan(&id, &userUUID, &createdAt, &lastSeenAt)
	session := model.Session{
		ID:        id,
		UserUUID:  userUUID,
		CreatedAt: createdAt,
	}
	if lastSeenAt.Valid {
		session.LastSeenAt = &lastSeenAt.Time
	}

	return session, err
}
