package sqldb

import (
	"context"
	"database/sql"
)

type txCtxKey struct{}

func withTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func txFromContext(ctx context.Context) (*sql.Tx, bool) {
	v := ctx.Value(txCtxKey{})
	if v == nil {
		return nil, false
	}
	tx, ok := v.(*sql.Tx)
	return tx, ok
}
