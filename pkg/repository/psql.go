package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mod/internal/config"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewClient(cntx context.Context, cfg *config.Config) (dbPool *pgxpool.Pool, err error) {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Db.Username, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Dbname)
	dbPool, err = pgxpool.Connect(cntx, dsn)
	if err != nil {
		return nil, err
	}
	return dbPool, nil
}
