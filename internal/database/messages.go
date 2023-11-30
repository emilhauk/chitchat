package database

import (
	"database/sql"
	"errors"
	"github.com/emilhauk/chitchat/internal/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type Messages struct {
	db *sql.DB

	create                        *sql.Stmt
	findForChannel                *sql.Stmt
	findLastMessageForChannelsSQL string
}

func NewMessageStore(db *sql.DB) Messages {
	create, err := db.Prepare("INSERT INTO messages (uuid, channel_uuid, user_uuid, content, version, sent_at) VALUE (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for messages.create")
	}

	findForChannel, err := db.Prepare("SELECT uuid, channel_uuid, user_uuid, content, version, sent_at, deleted_at, updated_at FROM messages WHERE channel_uuid = ? ORDER BY sent_at LIMIT ? OFFSET ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for messages.findForChannel")
	}

	// findLastMessageForChannelsSQL := "SELECT uuid, channel_uuid, user_uuid, content, version, MAX(sent_at), deleted_at, updated_at FROM messages WHERE channel_uuid IN (?) GROUP BY channel_uuid ORDER BY sent_at DESC"
	// TODO I expect this not to scale, but lets test it. Clever contributions are very welcome!
	findLastMessageForChannelsSQL := "SELECT m.* FROM messages m INNER JOIN (SELECT channel_uuid, MAX(sent_at) omg FROM messages GROUP BY channel_uuid) grouped_m ON m.channel_uuid=grouped_m.channel_uuid AND m.sent_at = grouped_m.omg AND m.channel_uuid IN (?)"
	_, err = db.Prepare(findLastMessageForChannelsSQL)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for messages.findLastMessageForChannels")
	}

	return Messages{
		db:                            db,
		create:                        create,
		findForChannel:                findForChannel,
		findLastMessageForChannelsSQL: findLastMessageForChannelsSQL,
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

func (s Messages) FindLastMessageForChannels(channelUUIDs ...string) ([]model.Message, error) {
	messages := make([]model.Message, 0)
	query, args, err := sqlx.In(s.findLastMessageForChannelsSQL, channelUUIDs)
	if err != nil {
		return messages, err
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return messages, err
	}

	for rows.Next() {
		message, err := s.mapToMessage(rows)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (s Messages) mapToMessage(row interface{ Scan(...any) error }) (model.Message, error) {
	var (
		uuid        string
		channelUUID string
		userUUID    string
		content     string
		version     uint32
		sentAt      time.Time
		deletedAt   sql.NullTime
		updatedAt   sql.NullTime
	)

	err := row.Scan(&uuid, &channelUUID, &userUUID, &content, &version, &sentAt, &deletedAt, &updatedAt)
	message := model.Message{
		UUID:        uuid,
		ChannelUUID: channelUUID,
		Sender:      model.User{UUID: userUUID},
		Content:     content,
		SentAt:      sentAt,
	}
	if deletedAt.Valid {
		message.DeletedAt = &deletedAt.Time
	}
	if updatedAt.Valid {
		message.UpdatedAt = &updatedAt.Time
	}

	return message, err
}
