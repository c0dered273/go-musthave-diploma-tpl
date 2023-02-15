package repositories

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	Save(ctx context.Context, u *models.User) error
	FindUserByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error)
}

type UsersRepositoryImpl struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UsersRepository {
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

func (r *UsersRepositoryImpl) FindUserByNameAndPasswd(ctx context.Context, name string, passwd string) (*models.User, error) {
	conn, err := getPgxConn(ctx, r)
	if err != nil {
		return nil, err
	}
	return store.FindUserByNameAndPasswd(ctx, conn, name, passwd)
}
