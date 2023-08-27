package store

import (
	"context"
	"strings"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	zerologAdapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

func NewPgxConn(ctx context.Context, logger zerolog.Logger, cfg *configs.ServerConfig) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(connStringWithOptions(cfg))
	if err != nil {
		return nil, err
	}

	pgxLogger := zerologAdapter.NewLogger(logger)
	pgxLogLevel, err := tracelog.LogLevelFromString(strings.ToLower(cfg.Database.LoggerLevel))
	if err != nil {
		return nil, err
	}

	tracer := &tracelog.TraceLog{
		Logger:   pgxLogger,
		LogLevel: pgxLogLevel,
	}
	connConfig.ConnConfig.Tracer = tracer

	poolConn, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	err = poolConn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return poolConn, nil
}

func connStringWithOptions(cfg *configs.ServerConfig) string {
	var sb strings.Builder
	isFirstParam := !strings.Contains(cfg.DatabaseURI, "?")
	sb.WriteString(cfg.DatabaseURI)

	if options := cfg.Database.Connection.Options; len(options) > 0 {
		for key, value := range options {
			if isFirstParam {
				sb.WriteRune('?')
				isFirstParam = false
			} else {
				sb.WriteRune('&')
			}

			sb.WriteString(key)
			sb.WriteRune('=')
			sb.WriteString(value)
		}
	}
	return sb.String()
}

type ConnCheck interface {
	Ping(ctx context.Context) error
}

type PgxConnCheck struct {
	Conn *pgxpool.Pool
}

func (p PgxConnCheck) Ping(ctx context.Context) error {
	return p.Conn.Ping(ctx)
}

func NewPgxConnCheck(conn *pgxpool.Pool) ConnCheck {
	return PgxConnCheck{Conn: conn}
}
