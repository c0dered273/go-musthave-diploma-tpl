package loggers

import (
	"os"
	"strings"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func NewDefaultLogger() zerolog.Logger {
	prettyOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(prettyOutput).With().Timestamp().Logger()
}

func NewServerLogger(cfg *configs.ServerConfig) zerolog.Logger {
	var logger zerolog.Logger

	if cfg.Server.Logger.Format == "pretty" {
		logger = NewDefaultLogger().With().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	if cfg.Server.Logger.Caller {
		logger = logger.With().Caller().Logger()
	}

	zerolog.SetGlobalLevel(stringToLevel(cfg.Server.Logger.Level))

	return logger.With().Str("appID", cfg.Server.Name).Logger()
}

func stringToLevel(lvl string) zerolog.Level {
	var level zerolog.Level
	switch strings.ToLower(lvl) {
	case "trace":
		level = zerolog.TraceLevel
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}
	return level
}
