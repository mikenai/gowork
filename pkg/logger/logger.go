package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level string `conf:"default:error"`
	Human bool
}

func New(cfg Config) (zerolog.Logger, error) {
	var output io.Writer = os.Stdout
	if cfg.Human {
		output = zerolog.ConsoleWriter{Out: os.Stderr}
	}
	log := zerolog.New(output)

	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return log, fmt.Errorf("failed to parse log: %w", err)
	}
	log = log.Level(lvl)

	return log, nil
}

func DefaultLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.ErrorLevel)
}
