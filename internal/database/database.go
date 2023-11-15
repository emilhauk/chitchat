package database

import (
	"database/sql"
	"fmt"
	"github.com/emilhauk/chitchat/config"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

var log = config.Logger

func NewConnectionPool(config config.DatabaseConfig) (*sql.DB, error) {
	dataSourceParams := url.Values{}
	dataSourceParams.Add("parseTime", "true")
	dataSourceParams.Add("time_zone", "'+00:00'")
	dataSourceParams.Add("loc", "UTC")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", config.Username, config.Password, config.Hostname, config.Port, config.Name, dataSourceParams.Encode())
	db, err := sql.Open(config.Driver, dataSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DBStore struct {
	Users       Users
	Credentials Credentials
	Sessions    Sessions
	Channels    Channels
	Messages    Messages
}

func NewDBStore(db *sql.DB) DBStore {
	return DBStore{
		Users:       NewUserStore(db),
		Credentials: NewCredentialStore(db),
		Sessions:    NewSessionStore(db),
		Channels:    NewChannelStore(db),
		Messages:    NewMessageStore(db),
	}
}
