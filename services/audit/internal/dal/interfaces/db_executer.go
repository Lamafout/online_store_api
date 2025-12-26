package interfaces

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DBExecuter interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
