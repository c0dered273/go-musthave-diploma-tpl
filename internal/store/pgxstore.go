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
	"github.com/shopspring/decimal"
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
	sql := `SELECT o.id, os.name, u.username, o.amount, o.uploaded_at
			FROM orders o
					 INNER JOIN users u on o.user_id = u.id
					 INNER JOIN order_status os on o.status_id = os.id
			WHERE o.id = $1;`

	var status string
	order := models.Order{}

	err := conn.QueryRow(ctx, strip(sql), orderID).Scan(&order.ID, &status, &order.Username, &order.Amount, &order.UploadedAt)
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

func GetUserBalance(ctx context.Context, conn *pgxpool.Pool, username string) (decimal.Decimal, error) {
	sql := "SELECT u.balance FROM users u WHERE u.username = $1"

	var balance decimal.Decimal
	err := conn.QueryRow(ctx, sql, username).Scan(&balance)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}

func FindOrdersByUsername(ctx context.Context, conn *pgxpool.Pool, username string) (models.Orders, error) {
	sql := `SELECT o.id, os.name, u.username, o.amount, o.uploaded_at 
			FROM orders o 
			    INNER JOIN order_status os on o.status_id = os.id 
			    INNER JOIN users u on o.user_id = u.id 
			WHERE u.username = $1
			ORDER BY o.uploaded_at`

	orders := make([]models.Order, 0)

	rows, err := conn.Query(ctx, strip(sql), username)
	if err != nil {
		return nil, err
	}

	var status string
	o := models.Order{}
	for rows.Next() {
		err := rows.Scan(&o.ID, &status, &o.Username, &o.Amount, &o.UploadedAt)
		if err != nil {
			return nil, err
		}

		if os, err := models.ParseStatus(status); err != nil {
			return nil, err
		} else {
			o.Status = os
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func FindWithdrawalsByUsername(ctx context.Context, conn *pgxpool.Pool, username string) (models.Withdrawals, error) {
	sql := `SELECT w.order_id, u.username, w.amount, w.processed_at 
			FROM withdrawals w 
			    INNER JOIN users u on u.id = w.user_id 
			WHERE u.username = $1`

	withdrawals := make([]models.Withdrawal, 0)

	rows, err := conn.Query(ctx, strip(sql), username)
	if err != nil {
		return nil, err
	}

	w := models.Withdrawal{}
	for rows.Next() {
		err := rows.Scan(&w.OrderID, &w.Username, &w.Amount, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}

func GetAllWithdrawalsSumByUsername(ctx context.Context, conn *pgxpool.Pool, username string) (decimal.Decimal, error) {
	sql := "SELECT COALESCE(SUM(w.amount), 0) FROM withdrawals w INNER JOIN users u on w.user_id = u.id WHERE u.username = $1"

	var allWithdrawal decimal.Decimal
	err := conn.QueryRow(ctx, sql, username).Scan(&allWithdrawal)
	if err != nil {
		return decimal.Zero, err
	}

	return allWithdrawal, nil
}
