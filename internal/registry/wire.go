package registry

import (
	"log/slog"
	"time"

	"go-clean-api/internal/app"
	"go-clean-api/internal/app/interceptor/logging"
	"go-clean-api/internal/app/interceptor/validation"
	"go-clean-api/internal/registry/interceptor"
	"go-clean-api/pkg/sqldb"
	"go-clean-api/pkg/validate"
)

type deps struct {
	dbClient  *sqldb.Client
	validator *validate.Validator
	logger    *slog.Logger
	nowFunc   func() time.Time
}

// applyStandard は読み取り系（Query）や Tx が不要なユースケースに、
// 標準の横断関心（ログ・バリデーション）だけを合成して返します。
//
// 実行順（外側 → 内側 → コア）: logging → validation → uc
func applyStandard[I, O any](name string, uc app.Interactor[I, O], d *deps) app.Interactor[I, O] {
	return applyInterceptors(uc,
		logging.Wrap[I, O](name, d.logger, d.nowFunc),
		validation.Wrap[I, O](d.validator),
	)
}

// applyStandardWithRequiredTx は登録・更新など DB 書き込みを伴うユースケース向けです。
// applyStandard に加え、コアの直前でトランザクションを張ります。
//
// ctx に既に tx があればそれを使い、なければ Begin します（RequiredTx）。
// 実行順: logging → validation → RequiredTx → uc
func applyStandardWithRequiredTx[I, O any](name string, uc app.Interactor[I, O], d *deps) app.Interactor[I, O] {
	uc = applyInterceptors(uc,
		interceptor.RequiredTx[I, O](d.dbClient),
	)
	return applyStandard(name, uc, d)
}

// applyStandardWithRequiredNewTx は applyStandardWithRequiredTx と同様ですが、
// ctx に tx があっても常に新しいトランザクションを開始します（RequiredNewTx）。
// ネストした処理を別トランザクションで切りたいときに使います。
/*func applyStandardWithRequiredNewTx[I, O any](name string, uc app.Interactor[I, O], d *deps) app.Interactor[I, O] {
	uc = applyInterceptors(uc,
		interceptor.RequiredNewTx[I, O](d.dbClient),
	)
	return applyStandard(name, uc, d)
}*/

// applyInterceptors は与えられた Wrap をコアに重ねて 1 本の Interactor にします。
// スライス末尾の interceptor がコアに最も近く（先に Execute されるのは外側）なります。
func applyInterceptors[I, O any](
	core app.Interactor[I, O],
	interceptors ...app.Wrap[I, O],
) app.Interactor[I, O] {
	instance := core
	for i := len(interceptors) - 1; i >= 0; i-- {
		instance = interceptors[i](instance)
	}
	return instance
}
