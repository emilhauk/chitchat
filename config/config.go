package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

// AppConfig is a struct for storing configuration for the application
type AppConfig struct {
	Port      int
	Host      string
	PublicURL string
	InDevMode bool
	Version   string
	Location  *time.Location
}

// DatabaseConfig stores configuration for a specific database connection
type DatabaseConfig struct {
	Driver   string
	Hostname string
	Name     string
	Username string
	Password string
	Port     int
}

type SMTPConfig struct {
	Enabled bool
	Host    string
	Port    string
}

var (
	// The App Contains configuration for this application
	App AppConfig

	// Database contains configuration for the database
	Database DatabaseConfig

	Mail SMTPConfig

	Logger zerolog.Logger

	version string
)

func init() {
	if version == "" {
		version = "unknown"
	}

	europeOslo, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		panic(err)
	}

	App = AppConfig{
		Host:      envString("HOSTNAME", "localhost"),
		Port:      envInt("PORT", 3333),
		Version:   version,
		Location:  europeOslo,
		InDevMode: envString("MODE", "development") == "development",
	}
	App.PublicURL = envString("PUBLIC_URL", fmt.Sprintf("http://%s:%d", App.Host, App.Port))

	Database = DatabaseConfig{
		Driver:   envString("DB_DRIVER", "mysql"),
		Hostname: envString("DB_HOSTNAME", "localhost"),
		Name:     envString("DB_NAME", "chitchat"),
		Username: envString("DB_USERNAME", "chitchat"),
		Password: envString("DB_PASSWORD", "password"),
		Port:     envInt("DB_PORT", 3390),
	}

	Mail = SMTPConfig{
		Enabled: envBool("SMTP_ENABLED", false),
		Host:    envString("SMTP_HOST", ""),
		Port:    envString("SMTP_PORT", "25"),
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if App.InDevMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		}).With().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		Logger = log.Logger
	}
}

func envString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func envBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func envInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		panic(err)
	}
	return value
}
