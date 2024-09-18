package eventing

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BeginTxFunc[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(tx pgx.Tx) (T, error),
) (T, error) {
	var tx pgx.Tx
	tx, err := pool.Begin(ctx)
	if err != nil {
		return *new(T), err
	}

	return beginFuncExec(ctx, tx, fn)
}

func beginFuncExec[T any](ctx context.Context, tx pgx.Tx, fn func(pgx.Tx) (T, error)) (ret T, err error) {
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = rollbackErr
		}
	}()

	ret, fErr := fn(tx)
	if fErr != nil {
		_ = tx.Rollback(ctx) // ignore rollback error as there is already an error to return
		return ret, fErr
	}

	return ret, tx.Commit(ctx)
}
