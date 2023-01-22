package controller

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5/pgconn"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
