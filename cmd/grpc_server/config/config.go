package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/mikenai/gowork/pkg/logger"
)

type Help string

func (h Help) String() string {
	return string(h)
}

type Config struct {
	Port string `conf:"default::5050"`

	Log logger.Config
	DB  SQLite
}

type SQLite struct {
	DSN string `conf:"default:tmp/db.sqlite3"`

	MaxIdleConns    int           `conf:"default:1"`
	MaxOpenConns    int           `conf:"default:5"`
	ConnMaxLifetime time.Duration `conf:"default:5s"`
	ConnMaxIdleTime time.Duration `conf:"default:5s"`
}

func New() (Config, Help, error) {
	cfg := Config{}

	if help, err := conf.Parse("", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return Config{}, Help(help), err
		}
		return Config{}, "", fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, "", nil
}
