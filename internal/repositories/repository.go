package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxTransactionKey struct{}

type PgxConn *pgxpool.Pool

type Repository interface {
	GetConn() *pgxpool.Pool
}

var ErrInvalidTxType = errors.New("repository: invalid tx type, tx type should be *pgxpool.Pool")

func getPgxConn(ctx context.Context, r Repository) (*pgxpool.Pool, error) {
	txFromContext := ctx.Value(ctxTransactionKey{})
	if txFromContext == nil {
		return r.GetConn(), nil
	}
	if tx, ok := txFromContext.(*pgxpool.Pool); ok {
		return tx, nil
	}
	return nil, ErrInvalidTxType
}

func withSQLTransaction(ctx context.Context, r Repository, fn func(ctx context.Context) error) error {
	tx, err := r.GetConn().Begin(ctx)
	if err != nil {
		return err
	}
	trxCtx := context.WithValue(ctx, ctxTransactionKey{}, tx)

	err = fn(trxCtx)
	if err != nil {
		return tx.Rollback(trxCtx)
	}
	return tx.Commit(trxCtx)
}
