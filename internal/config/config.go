package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// EventJournalPGConfig configures the main event journal postgres instance.
type EventJournalPGConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

// ConnStr returns the formatted postgres connection string.
func (c *EventJournalPGConfig) ConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s", c.User, c.Password, c.Host, c.Name)
}

// EventJournalRedisConfig configures the main redis instance.
type EventJournalRedisConfig struct {
	Host string
}

// EventJournalConfig configures the event journal.
type EventJournalConfig struct {
	PG    EventJournalPGConfig
	Redis EventJournalRedisConfig
}

type PermifyConfig struct {
	Host string
}

// SetupConfig configures the server setup.
type SetupConfig struct {
	Root SetupRootConfig
}

// SetupRootConfig configures the server root user.
type SetupRootConfig struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// Config configures the server.
type Config struct {
	EventJournal EventJournalConfig
	Permify      PermifyConfig
	Setup        SetupConfig

	Host string
}

func (c *Config) Validate() error {
	if c.EventJournal.PG.Host == "" {
		return fmt.Errorf("EventJournal.PG.Host is required")
	}
	if c.EventJournal.PG.User == "" {
		return fmt.Errorf("EventJournal.PG.User is required")
	}
	if c.EventJournal.PG.Password == "" {
		return fmt.Errorf("EventJournal.PG.Password is required")
	}
	if c.EventJournal.PG.Name == "" {
		return fmt.Errorf("EventJournal.PG.Name is required")
	}

	if c.EventJournal.Redis.Host == "" {
		return fmt.Errorf("EventJournal.Redis.Host is required")
	}

	if c.Permify.Host == "" {
		return fmt.Errorf("Permify.Host is required")
	}

	if c.Setup.Root.Email == "" {
		return fmt.Errorf("Setup.Root.Email is required")
	}
	return nil
}

func NewConfig(viper *viper.Viper) (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &config, nil
}
