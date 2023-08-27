package repositories

import (
	"context"
	"errors"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type OrderRepository interface {
	Save(ctx context.Context, order *models.Order) error
	FindByID(ctx context.Context, orderID uint64) (*models.Order, error)
	FindByUsername(ctx context.Context, username string) (models.Orders, error)
	UpdateByID(ctx context.Context, orderID uint64, status models.OrderStatus, amount decimal.Decimal) error
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

func (r *OrderRepositoryImpl) Save(ctx context.Context, order *models.Order) error {
	sql := `INSERT INTO orders(id, status_id, user_id, amount, uploaded_at) 
			VALUES ($1,
			        (SELECT os.id FROM order_status os WHERE os.name = $2),
			        (SELECT u.id FROM users u WHERE u.username = $3),
			        $4,
			        $5)
			ON CONFLICT DO NOTHING`

	commandTag, err := r.Conn.Exec(ctx, strip(sql), order.ID, order.Status, order.Username, order.Amount, order.UploadedAt)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrAlreadyExists
	}

	return nil
}

func (r *OrderRepositoryImpl) FindByID(ctx context.Context, orderID uint64) (*models.Order, error) {
	sql := `SELECT o.id, os.name, u.username, o.amount, o.uploaded_at
			FROM orders o
					 INNER JOIN users u on o.user_id = u.id
					 INNER JOIN order_status os on o.status_id = os.id
			WHERE o.id = $1;`

	var status string
	order := models.Order{}

	err := r.Conn.QueryRow(ctx, strip(sql), orderID).Scan(&order.ID, &status, &order.Username, &order.Amount, &order.UploadedAt)
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

func (r *OrderRepositoryImpl) FindByUsername(ctx context.Context, username string) (models.Orders, error) {
	sql := `SELECT o.id, os.name, u.username, o.amount, o.uploaded_at 
			FROM orders o 
			    INNER JOIN order_status os on o.status_id = os.id 
			    INNER JOIN users u on o.user_id = u.id 
			WHERE u.username = $1
			ORDER BY o.uploaded_at`

	orders := make([]models.Order, 0)

	rows, err := r.Conn.Query(ctx, strip(sql), username)
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

func (r *OrderRepositoryImpl) UpdateByID(ctx context.Context, orderID uint64, status models.OrderStatus, amount decimal.Decimal) error {
	sql := `UPDATE orders SET
                  status_id = (SELECT os.id FROM order_status os WHERE os.name = $2),
                  amount = $3
              WHERE ID = $1`

	tag, err := r.Conn.Exec(ctx, sql, orderID, status, amount)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
