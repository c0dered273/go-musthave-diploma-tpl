package repositories

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WithdrawalRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	FindByUsername(ctx context.Context, username string) (models.Withdrawals, error)
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

func (r *WithdrawalRepositoryImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withSQLTransaction(ctx, r, fn)
}

//func (r *WithdrawalRepositoryImpl) Save(ctx context.Context, w *models.Withdrawal) error {
//	conn, err := getPgxConn(ctx, r)
//	if err != nil {
//		return err
//	}
//	return store.SaveOrder(ctx, conn, order)
//}
//

func (r *WithdrawalRepositoryImpl) FindByUsername(ctx context.Context, username string) (models.Withdrawals, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return nil, err
	}
	return store.FindWithdrawalsByUsername(ctx, conn, username)
}
