package transactor

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

type txkey string

const key = txkey("tx")

func (tm *TransactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
	if err != nil {
		return err
	}

	if err = fx(context.WithValue(ctx, key, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err = tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}
