// Package service は複数の usecase/query を組み合わせる Facade 層を提供します。
//
// この層は単一の usecase では表現しづらい複合処理を担当し、
// ユースケース間のオーケストレーションを行います。
//
// 現状:
//   - このディレクトリには doc のみ（実装は未配置）。必要になったらサブパッケージを追加する。
//
// ルール:
//   - app.Interactor[I, O] を実装する形で Facade を置く。
//   - 命名は Input/Output を使う（Request/Response は delivery 層で使用）。
//   - 単純な処理は usecase/query に置き、複合ロジックのみを service に置く。
//   - インフラ依存の処理は直接持たず、必要な依存は registry で注入する。
package service
