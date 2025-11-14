package logger

import (
	"github.com/rs/zerolog"
)

// WithName creates sub-logger and setting name attribute
func WithName(logger zerolog.Logger, name string) zerolog.Logger {
	return logger.With().Str("name", name).Logger()
}
