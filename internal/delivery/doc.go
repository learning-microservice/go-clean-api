// Package delivery は外部クライアント向けの入出力（transport）を担当します。
//
// この層は HTTP ハンドラ・ルーティング・リクエスト/レスポンスの変換を行い、
// アプリケーションロジック（usecase / query）の呼び出しに専念します。
//
// 構成:
//   - restapi: Echo による REST API。バージョンごとにサブパッケージを切る（例: v1/auth）。
//   - API 契約の型は api/openapi（oapi-codegen 生成）を利用する。
//
// 例:
//   - 認証: internal/delivery/restapi/v1/auth（Login / Register）が usecase/auth を呼び出す。
//
// ルール:
//   - ビジネスルール・永続化ロジックは持たない（変換と HTTP ステータス・エラー応答のマッピングのみ）。
//   - 命名は Request/Response（OpenAPI 型）を delivery で、Input/Output は usecase/query で使い分ける。
//   - ドメインエラー（domain/errors）は HTTP ステータスと ErrorResponse に変換して返す。
//   - usecase / query の Execute には c.Request().Context() を渡す。
//   - domain や infra を直接参照しない（registry 経由で注入された Interactor のみ呼ぶ）。
package delivery
