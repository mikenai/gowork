package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/mikenai/gowork/common/logger"
)

type Help string

func (h Help) String() string {
	return string(h)
}

type Config struct {
	GracefullTimeout time.Duration `conf:"default:30s"`

	Log  logger.Config
	HTTP HTTP
}

type HTTP struct {
	Addr         string        `conf:"default::8181"`
	ReadTimeout  time.Duration `conf:"default:1s"`
	WriteTimeout time.Duration `conf:"default:1s"`
	IdleTimeout  time.Duration `conf:"default:5s"`
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
