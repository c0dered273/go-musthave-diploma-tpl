package repositories

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type WithdrawalRepository interface {
	FindByUsername(ctx context.Context, username string) (models.Withdrawals, error)
	GetAllWithdrawalByUsername(ctx context.Context, username string) (decimal.Decimal, error)
	Save(ctx context.Context, w *models.Withdrawal) error
}

type WithdrawalRepositoryImpl struct {
	Conn *pgxpool.Pool
}

func NewWithdrawalRepository(conn *pgxpool.Pool) WithdrawalRepository {
	return &WithdrawalRepositoryImpl{
		Conn: conn,
	}
}

func (r *WithdrawalRepositoryImpl) GetConn() *pgxpool.Pool {
	return r.Conn
}

func (r *WithdrawalRepositoryImpl) Save(ctx context.Context, w *models.Withdrawal) error {
	sql := `INSERT INTO withdrawals(order_id, user_id, amount, processed_at) 
				VALUES($1,
				       (SELECT u.id FROM users u WHERE u.username = $2),
				       $3,
				       $4)
				ON CONFLICT DO NOTHING`

	n, err := r.Conn.Exec(ctx, strip(sql), w.OrderID, w.Username, w.Amount, w.ProcessedAt)
	if err != nil {
		return err
	}

	if n.RowsAffected() == 0 {
		return ErrAlreadyExists
	}

	return nil
}

func (r *WithdrawalRepositoryImpl) FindByUsername(ctx context.Context, username string) (models.Withdrawals, error) {
	sql := `SELECT w.order_id, u.username, w.amount, w.processed_at 
			FROM withdrawals w 
			    INNER JOIN users u on u.id = w.user_id 
			WHERE u.username = $1`

	withdrawals := make([]models.Withdrawal, 0)

	rows, err := r.Conn.Query(ctx, strip(sql), username)
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

func (r *WithdrawalRepositoryImpl) GetAllWithdrawalByUsername(ctx context.Context, username string) (decimal.Decimal, error) {
	sql := "SELECT COALESCE(SUM(w.amount), 0) FROM withdrawals w INNER JOIN users u on w.user_id = u.id WHERE u.username = $1"

	var allWithdrawal decimal.Decimal
	err := r.Conn.QueryRow(ctx, sql, username).Scan(&allWithdrawal)
	if err != nil {
		return decimal.Zero, err
	}

	return allWithdrawal, nil
}
