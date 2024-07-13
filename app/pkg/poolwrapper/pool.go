package poolwrapper

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(pool *pgxpool.Pool) *Pool {
	p := &Pool{
		pool: pool,
	}

	return p
}

func (p Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Begin(ctx)
	}

	return p.pool.Begin(ctx)
}

func (p Pool) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
	}

	return p.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (p Pool) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	if tx := extractTx(ctx); tx != nil {
		return tx.SendBatch(ctx, b)
	}

	return p.pool.SendBatch(ctx, b)
}

func (p Pool) Exec(ctx context.Context, sql string, args ...any) (commandTag pgconn.CommandTag, err error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}

	return p.pool.Exec(ctx, sql, args...)
}

func (p Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}

	return p.pool.Query(ctx, sql, args...)
}

func (p Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}

	return p.pool.QueryRow(ctx, sql, args...)
}

func (p Pool) TxFunc(ctx context.Context, f func(context.Context, pgx.Tx) error) (err error) {
	if tx := extractTx(ctx); tx != nil {
		nestedTx, err := tx.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin nested transaction: %w", err)
		}

		defer func() {
			if rErr := nestedTx.Rollback(ctx); rErr != nil && !errors.Is(rErr, pgx.ErrTxClosed) {
				err = errors.Join(err, fmt.Errorf("rollback nested tx: %w", rErr))
			}
		}()

		err = f(injectTx(ctx, nestedTx), nestedTx)
		if err == nil {
			if err = nestedTx.Commit(ctx); err != nil {
				return fmt.Errorf("commit nested tx: %w", err)
			}
		}

		return err
	} else {
		tx, err = p.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}

		defer func() {
			if rErr := tx.Rollback(ctx); rErr != nil && !errors.Is(rErr, pgx.ErrTxClosed) {
				err = errors.Join(err, fmt.Errorf("rollback tx: %w", rErr))
			}
		}()

		err = f(injectTx(ctx, tx), tx)
		if err == nil {
			if err = tx.Commit(ctx); err != nil {
				return fmt.Errorf("commit tx: %w", err)
			}
		}

		return err
	}

}
