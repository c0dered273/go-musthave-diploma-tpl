package repositories

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	Save(ctx context.Context, order *models.Order) error
	FindByID(ctx context.Context, orderID uint64) (*models.Order, error)
	FindByUsername(ctx context.Context, username string) (models.Orders, error)
}

type OrderRepositoryImpl struct {
	Conn *pgxpool.Pool
}

func NewOrderRepository(conn *pgxpool.Pool) OrderRepository {
	return &OrderRepositoryImpl{
		Conn: conn,
	}
}

func (r *OrderRepositoryImpl) GetConn() *pgxpool.Pool {
	return r.Conn
}

func (r *OrderRepositoryImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withSQLTransaction(ctx, r, fn)
}

func (r *OrderRepositoryImpl) Save(ctx context.Context, order *models.Order) error {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return err
	}
	return store.SaveOrder(ctx, conn, order)
}

func (r *OrderRepositoryImpl) FindByID(ctx context.Context, orderID uint64) (*models.Order, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return nil, err
	}
	return store.FindOrderByID(ctx, conn, orderID)
}

func (r *OrderRepositoryImpl) FindByUsername(ctx context.Context, username string) (models.Orders, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return nil, err
	}
	return store.FindOrdersByUsername(ctx, conn, username)
}
