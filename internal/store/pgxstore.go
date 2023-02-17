package store

import (
	"context"
	"errors"
	"strings"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	zerologAdapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrNotFound           = errors.New("not found")
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
	isFirstParam := !strings.Contains(cfg.DatabaseUri, "?")
	sb.WriteString(cfg.DatabaseUri)

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

func strip(str string) string {
	rsl := strings.ReplaceAll(str, "\t", "")
	return rsl
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

func SaveUser(ctx context.Context, conn *pgxpool.Pool, u *models.User) error {
	sql := `INSERT INTO users(username, password, balance) 
				VALUES($1, crypt($2, gen_salt('bf')), 0) 
				ON CONFLICT DO NOTHING`

	n, err := conn.Exec(ctx, strip(sql), u.Username, u.Password)
	if err != nil {
		return err
	}

	if n.RowsAffected() == 0 {
		return ErrUserAlreadyExists
	}

	return nil
}

func FindUserByNameAndPasswd(
	ctx context.Context, conn *pgxpool.Pool, name string, passwd string,
) (*models.User, error) {
	sql := "SELECT username FROM users WHERE username=$1 AND password=crypt($2, password)"

	var username string
	err := conn.QueryRow(ctx, sql, name, passwd).Scan(&username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &models.User{
		Username: username,
		Password: "",
	}, nil
}

func SaveOrder(
	ctx context.Context,
	conn *pgxpool.Pool,
	order *models.Order,
) error {
	sql := `INSERT INTO orders(id, status_id, user_id, amount, uploaded_at) 
			VALUES ($1,
			        (SELECT os.id FROM order_status os WHERE os.name = $2),
			        (SELECT u.id FROM users u WHERE u.username = $3),
			        $4,
			        $5)
			ON CONFLICT DO NOTHING`

	commandTag, err := conn.Exec(ctx, strip(sql), order.ID, order.Status, order.Username, order.Amount, order.UploadedAt)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrOrderAlreadyExists
	}

	return nil
}

func FindOrderByID(ctx context.Context, conn *pgxpool.Pool, orderID uint64) (*models.Order, error) {
	sql := `SELECT o.id, u.username, os.name, o.amount, o.uploaded_at
			FROM orders o
					 INNER JOIN users u on o.user_id = u.id
					 INNER JOIN order_status os on o.status_id = os.id
			WHERE o.id = $1;`

	var status string
	order := models.Order{}

	err := conn.QueryRow(ctx, strip(sql), orderID).Scan(&order.ID, &order.Username, &status, &order.Amount, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if os, err := models.ParseStatus(status); err != nil {
		return nil, err
	} else {
		order.Status = os
	}

	return &order, nil
}
