// Package usecase は書き込み系（Command）のアプリケーションロジックを提供します。
//
// この層はドメインルールに基づく状態変更を担当し、更新系ユースケースを実装します。
// 読み取り最適化された検索処理は query 層で扱います。
//
// サブパッケージの切り方:
//   - ドメインの集約パッケージ名と 1:1 である必要はない。
//   - 機能・ユースケース単位でまとめる（例: auth はログインと登録を扱い、domain/user を利用）。
//     domain/user とパッケージ名がぶつからないよう、import の衝突にも注意する。
//
// ルール:
//   - Usecase の公開 API は Execute(ctx, ...) を中心にし、registry からは必要に応じて
//     app.Interactor としてラップして合成する。
//   - 公開ユースケース型は auth と同様、次のように型エイリアスで app.Interactor と同義にできる。
//     type FooUsecase = app.Interactor[*FooInput, *FooOutput]
//     mockgen -source はインターフェース宣言のみ対象のため、エイリアスだけのユースケース型は Mock が
//     生成されない。テストでは Output Port や domain.Repository を mock するか、delivery 向けに
//     手書き fake（Execute のみ）を使う。
//   - 命名は Input/Output を使う（Request/Response は delivery 層で使用）。
//   - DB の実装詳細や transport 型は持ち込まない（ドメインの Repository などポートのみ）。
//   - 横断関心事は internal/registry/wire.go で合成する（applyStandard /
//     applyStandardWithRequiredTx）。logging / validation は internal/app/interceptor、
//     トランザクションは internal/registry/interceptor。
package usecase
