package repositories

import (
	"context"
	"errors"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var (
	ErrBalanceNotEnough = errors.New("user balance not enough")
)

type UserRepository interface {
	Save(ctx context.Context, u *models.User) error
	FindByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error)
	GetBalance(ctx context.Context, username string) (decimal.Decimal, error)
	Withdrawing(ctx context.Context, username string, orderID string, amount decimal.Decimal) error
}

type UsersRepositoryImpl struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &UsersRepositoryImpl{
		Conn: conn,
	}
}

func (r *UsersRepositoryImpl) GetConn() *pgxpool.Pool {
	return r.Conn
}

func (r *UsersRepositoryImpl) Save(ctx context.Context, u *models.User) error {
	sql := `INSERT INTO users(username, password, balance) 
				VALUES($1, crypt($2, gen_salt('bf')), 0) 
				ON CONFLICT DO NOTHING`

	n, err := r.Conn.Exec(ctx, strip(sql), u.Username, u.Password)
	if err != nil {
		return err
	}

	if n.RowsAffected() == 0 {
		return ErrAlreadyExists
	}

	return nil
}

func (r *UsersRepositoryImpl) FindByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error) {
	sql := "SELECT username FROM users WHERE username=$1 AND password=crypt($2, password)"

	var username string
	err := r.Conn.QueryRow(ctx, sql, name, passwd).Scan(&username)
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

func (r *UsersRepositoryImpl) GetBalance(ctx context.Context, username string) (decimal.Decimal, error) {
	sql := "SELECT u.balance FROM users u WHERE u.username = $1"

	var balance decimal.Decimal
	err := r.Conn.QueryRow(ctx, sql, username).Scan(&balance)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}

func (r *UsersRepositoryImpl) Withdrawing(ctx context.Context, username string, orderID string, amount decimal.Decimal) error {
	sql := "CALL withdraw_from_user_balance($1, $2, $3)"

	_, err := r.Conn.Exec(ctx, sql, username, orderID, amount)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Message == "balance not enough" {
			return ErrBalanceNotEnough
		}
		return err
	}

	return nil
}
