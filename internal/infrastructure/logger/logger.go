package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func New() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &logger
}
