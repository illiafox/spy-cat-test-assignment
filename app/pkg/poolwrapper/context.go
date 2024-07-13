package poolwrapper

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type connKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, connKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	tx := ctx.Value(connKey{})
	if tx != nil {
		return tx.(pgx.Tx)
	}

	return nil
}
