// Package interceptor は app.Interactor に横断関心事を合成する仕組みを提供します。
//
// この層は delivery に依存せず、ユースケース実行の前後に共通処理を差し込みます。
// internal/registry/interceptor パッケージ（別物）が sqldb 等を使う合成を担当する。
//
// 提供する interceptor（registry/wire.go の applyStandard* から利用）:
//   - logging: 実行時間・成功/失敗のログ
//   - validation: 入力 struct の検証（失敗時は検証エラーを返す）
//
// ルール:
//   - app.Wrap[I, O] と app.Interactor[I, O] を組み合わせて合成する。
//   - DB の実装詳細・*sql.Tx・HTTP の型は持ち込まない。
package interceptor
