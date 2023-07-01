package config

import (
	"fmt"

	"github.com/ardanlabs/conf/v3"
)

type Config struct {
	Port string `conf:"default:8081"`
}

func LoadConfig() (c Config, err error) {
	if _, err = conf.Parse("", &c); err != nil {
		return c, fmt.Errorf("failed to parse config: %w", err)
	}
	return
}
