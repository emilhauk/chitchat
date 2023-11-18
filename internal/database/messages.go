package database

import (
	"database/sql"
	"errors"
	"github.com/emilhauk/chitchat/internal/model"
	"time"
)

type Messages struct {
	db *sql.DB

	create         *sql.Stmt
	findForChannel *sql.Stmt
}

func NewMessageStore(db *sql.DB) Messages {
	create, err := db.Prepare("INSERT INTO messages (uuid, channel_uuid, user_uuid, content, version, sent_at) VALUE (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for messages.create")
	}

	findForChannel, err := db.Prepare("SELECT uuid, user_uuid, content, version, sent_at, deleted_at, updated_at FROM messages WHERE channel_uuid = ? ORDER BY sent_at LIMIT ? OFFSET ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for messages.findForChannel")
	}

	return Messages{
		db:             db,
		create:         create,
		findForChannel: findForChannel,
	}
}

func (s Messages) Create(channelUUID string, m model.Message) error {
	_, err := s.create.Exec(m.UUID, channelUUID, m.Sender.UUID, m.Content, m.Version, m.SentAt)
	return err
}

func (s Messages) FindForChannel(channelUUID string, limit, offset int32) ([]model.Message, error) {
	messages := make([]model.Message, 0)
	rows, err := s.findForChannel.Query(channelUUID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages, nil
		}
		return messages, err
	}
	for rows.Next() {
		channel, err := s.mapToMessage(rows)
		if err != nil {
			return messages, err
		}
		messages = append(messages, channel)
	}
	return messages, nil
}

func (s Messages) mapToMessage(row interface{ Scan(...any) error }) (model.Message, error) {
	var (
		uuid      string
		userUUID  string
		content   string
		version   uint32
		sentAt    time.Time
		deletedAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := row.Scan(&uuid, &userUUID, &content, &version, &sentAt, &deletedAt, &updatedAt)
	message := model.Message{
		UUID:    uuid,
		Sender:  model.User{UUID: userUUID},
		Content: content,
		SentAt:  sentAt,
	}
	if deletedAt.Valid {
		message.DeletedAt = &deletedAt.Time
	}
	if updatedAt.Valid {
		message.UpdatedAt = &updatedAt.Time
	}

	return message, err
}
