package repositories

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type UserRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	Save(ctx context.Context, u *models.User) error
	FindByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error)
	GetUserBalance(ctx context.Context, username string) (decimal.Decimal, error)
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

func (r *UsersRepositoryImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return withSQLTransaction(ctx, r, fn)
}

func (r *UsersRepositoryImpl) Save(ctx context.Context, u *models.User) error {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return err
	}
	return store.SaveUser(ctx, conn, u)
}

func (r *UsersRepositoryImpl) FindByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return nil, err
	}
	return store.FindUserByNameAndPasswd(ctx, conn, name, passwd)
}

func (r *UsersRepositoryImpl) GetUserBalance(ctx context.Context, username string) (decimal.Decimal, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return decimal.Zero, err
	}
	return store.GetUserBalance(ctx, conn, username)
}
