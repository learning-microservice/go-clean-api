package interceptor

import (
	"context"

	"go-clean-api/internal/app"
	"go-clean-api/pkg/sqldb"
)

// RequiredTx: ctx に tx があればそのまま next。なければ Begin → WithTx → Commit/Rollback
func RequiredTx[I, O any](client *sqldb.Client) app.Wrap[I, O] {
	return func(next app.Interactor[I, O]) app.Interactor[I, O] {
		return transaction[I, O]{
			client: client,
			next:   next,
			always: false,
		}
	}
}

// RequiredNewTx: 常に db.BeginTx。ctx に既に tx があっても無視して新規（別トランザクション）
func RequiredNewTx[I, O any](client *sqldb.Client) app.Wrap[I, O] {
	return func(next app.Interactor[I, O]) app.Interactor[I, O] {
		return transaction[I, O]{
			client: client,
			next:   next,
			always: true,
		}
	}
}

type transaction[I, O any] struct {
	client *sqldb.Client
	next   app.Interactor[I, O]
	always bool
}

func (t transaction[I, O]) Execute(ctx context.Context, in I) (o O, err error) {
	// 既にトランザクションに参加している場合はオーナーシップを持たないが、
	// always が true の場合にのみ常に新規トランザクションを開始
	txFunc := t.client.RequiredTx
	if t.always {
		txFunc = t.client.RequiredNewTx
	}
	err = txFunc(ctx, func(txCtx context.Context) error {
		var e error
		o, e = t.next.Execute(txCtx, in)
		return e
	})
	return o, err
}
