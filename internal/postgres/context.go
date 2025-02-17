package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

// TransactionKey is the key used to store a pgx.Tx in a context.Context.
var TransactionKey = txKey{}

// WithTx returns a new context with the provided transaction.
// Use this function when you start a transaction.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TransactionKey, tx)
}

// DB is an interface that both pgx.Tx and *pgxpool.Pool satisfy.
// This allows our methods to work transparently with either a transaction or the pool.
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

// GetDBFromContext returns a DB instance that is either the transaction stored in the context or the provided pool.
func GetDBFromContext(ctx context.Context, pool *pgxpool.Pool) DB {
	if tx, ok := ctx.Value(TransactionKey).(pgx.Tx); ok {
		return tx
	}
	return pool
}
