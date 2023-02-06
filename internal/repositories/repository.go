package repositories

import (
	"context"
	"strings"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	Close()
	Ping(ctx context.Context) error
}

type CrudRepository struct {
	Conn   *pgxpool.Pool
	Logger zerolog.Logger
}

func NewCrudRepository(ctx context.Context, logger zerolog.Logger, cfg *configs.ServerConfig) (Repository, error) {
	var sb strings.Builder
	sb.WriteString(cfg.DatabaseUri)
	isFirstParam := !strings.Contains(sb.String(), "?")

	if options := cfg.Database.Connection.PgxOptions; len(options) > 0 {
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

	connConf, err := pgxpool.ParseConfig(sb.String())
	if err != nil {
		return nil, err
	}
	connPool, err := pgxpool.NewWithConfig(ctx, connConf)
	if err != nil {
		return nil, err
	}
	err = connPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &CrudRepository{
		Conn:   connPool,
		Logger: logger,
	}, nil
}

func (c *CrudRepository) Close() {
	c.Conn.Close()
}

func (c *CrudRepository) Ping(ctx context.Context) error {
	err := c.Conn.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
