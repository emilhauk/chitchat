package database

import (
	"database/sql"
	"errors"
	"github.com/emilhauk/chitchat/internal"
	"github.com/emilhauk/chitchat/internal/model"
	"time"
)

type Channels struct {
	db *sql.DB

	create         *sql.Stmt
	findByUUID     *sql.Stmt
	findForUser    *sql.Stmt
	findAllForUser *sql.Stmt

	addMember  *sql.Stmt
	findMember *sql.Stmt
}

func NewChannelStore(db *sql.DB) Channels {
	create, err := db.Prepare("INSERT INTO channels (uuid, name, created_at) VALUE (?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channels.create")
	}
	findByUUID, err := db.Prepare("SELECT uuid, name, created_at, updated_at FROM channels WHERE uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channels.findByUUID")
	}
	findForUser, err := db.Prepare("SELECT c.uuid, c.name, c.created_at, c.updated_at FROM channels c INNER JOIN channel_members cm ON c.uuid = cm.channel_uuid WHERE c.uuid = ? AND cm.user_uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channels.findForUser")
	}
	findAllForUser, err := db.Prepare("SELECT c.uuid, c.name, c.created_at, c.updated_at FROM channels c INNER JOIN channel_members cm ON c.uuid = cm.channel_uuid WHERE cm.user_uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channels.findAllForUser")
	}

	addMember, err := db.Prepare("INSERT INTO channel_members (channel_uuid, user_uuid, role, created_at) VALUE (?, ?, ?, ?)")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channel_members.addMember")
	}
	findMember, err := db.Prepare("SELECT * FROM channel_members WHERE channel_uuid = ? AND user_uuid = ?")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to prepare statement for channel_members.findMember")
	}

	return Channels{
		db:             db,
		create:         create,
		findByUUID:     findByUUID,
		findForUser:    findForUser,
		findAllForUser: findAllForUser,
		addMember:      addMember,
		findMember:     findMember,
	}
}

func (s Channels) Create(m model.Channel) error {
	_, err := s.create.Exec(m.UUID, m.Name, m.CreatedAt)
	return err
}

func (s Channels) FindByUUID(uuid string) (model.Channel, error) {
	channel, err := s.mapToChannel(s.findByUUID.QueryRow(uuid))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return channel, app.ErrChannelNotFound
	}
	return channel, err
}

func (s Channels) FindForUser(channelUUID, userUUID string) (model.Channel, error) {
	channel, err := s.mapToChannel(s.findForUser.QueryRow(channelUUID, userUUID))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return channel, app.ErrChannelNotFound
	}
	return channel, err
}

func (s Channels) FindAllForUser(userUUID string) ([]model.Channel, error) {
	channels := make([]model.Channel, 0)
	rows, err := s.findAllForUser.Query(userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return channels, nil
		}
		return channels, err
	}
	for rows.Next() {
		channel, err := s.mapToChannel(rows)
		if err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (s Channels) AddMember(channel model.Channel, user model.User, role model.ChannelRole) error {
	_, err := s.addMember.Exec(channel.UUID, user.UUID, role, time.Now())
	return err
}

func (s Channels) FindMember(channelUUID, userUUID string) (model.Member, error) {
	member, err := s.mapToMember(s.findMember.QueryRow(channelUUID, userUUID))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return member, app.ErrMemberNotFound
	}
	return member, err
}

func (s Channels) mapToChannel(row interface{ Scan(...any) error }) (model.Channel, error) {
	var (
		uuid      string
		name      sql.NullString
		createdAt time.Time
		updatedAt sql.NullTime
	)

	err := row.Scan(&uuid, &name, &createdAt, &updatedAt)
	channel := model.Channel{
		UUID:      uuid,
		CreatedAt: createdAt,
	}
	if name.Valid && name.String != "" {
		channel.Name = name.String
	}
	if updatedAt.Valid {
		channel.UpdatedAt = &updatedAt.Time
	}
	return channel, err
}

func (s Channels) mapToMember(row interface{ Scan(...any) error }) (model.Member, error) {
	var (
		channelUUID string
		userUUID    string
		role        model.ChannelRole
		createdAt   time.Time
		updatedAt   sql.NullTime
	)

	err := row.Scan(&channelUUID, &userUUID, &role, &createdAt, &updatedAt)

	member := model.Member{
		ChannelUUID: channelUUID,
		UserUUID:    userUUID,
		Role:        role,
		CreatedAt:   createdAt,
	}
	if updatedAt.Valid {
		member.UpdatedAt = &updatedAt.Time
	}
	return member, err
}
