package logger

import "github.com/rs/zerolog"

type Simple struct {
	logger zerolog.Logger
}

func NewSimple(logger zerolog.Logger) *Simple {
	return &Simple{logger: logger}
}

func (l Simple) Logf(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}
