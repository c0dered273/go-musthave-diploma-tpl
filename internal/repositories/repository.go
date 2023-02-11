package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/rs/zerolog"
)

type Repository interface {
	SqlxDB() *sqlx.DB
	Close() error
	Ping(ctx context.Context) error
}

type CrudRepository struct {
	DB     *sqlx.DB
	Logger zerolog.Logger
}

func (c *CrudRepository) Close() error {
	err := c.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *CrudRepository) Ping(ctx context.Context) error {
	err := c.DB.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (c *CrudRepository) SqlxDB() *sqlx.DB {
	return c.DB
}

type Hooks struct {
	logger zerolog.Logger
}

func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "begin", time.Now()), nil
}

func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)

	h.logger.Debug().
		Str("query", query).
		Str("args", fmt.Sprintf("%q", args)).
		Str("duration", fmt.Sprintf("%s", time.Since(begin))).
		Send()

	return ctx, nil
}

func NewCrudRepository(ctx context.Context, logger zerolog.Logger, cfg *configs.ServerConfig) (Repository, error) {
	sqlLogger := logger.With().Str("service", "sql").Logger()
	sql.Register("pgxWithHooks", sqlhooks.Wrap(&stdlib.Driver{}, &Hooks{logger: sqlLogger}))

	sqlxDB, err := sqlx.ConnectContext(ctx, "pgx", cfg.DatabaseUri)
	if err != nil {
		return nil, err
	}

	maxConns, ok := cfg.Database.Connection.Options["max_open_conns"]
	if ok {
		i, cErr := strconv.Atoi(maxConns)
		if cErr != nil {
			return nil, cErr
		}
		sqlxDB.SetMaxOpenConns(i)
	}
	maxLifeTime, ok := cfg.Database.Connection.Options["max_conn_lifetime_time"]
	if ok {
		d, cErr := time.ParseDuration(maxLifeTime)
		if cErr != nil {
			return nil, cErr
		}
		sqlxDB.SetConnMaxLifetime(d)
	}

	return &CrudRepository{
		DB:     sqlxDB,
		Logger: sqlLogger,
	}, nil
}
