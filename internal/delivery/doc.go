// Package delivery は外部クライアント向けの入出力（transport）を担当します。
//
// この層は HTTP ハンドラ・ルーティング・リクエスト/レスポンスの変換を行い、
// アプリケーションロジック（usecase / query）の呼び出しに専念します。
//
// 構成:
//   - restapi: Echo による REST API。バージョンごとにサブパッケージを切る（例: v1/auth）。
//   - restapi/httperror: エラー応答の JSON 化（Encode）とアクセスログ用属性（LogAttrs）。
//   - restapi/middleware/accesslog: RequestLogger（HandleError: true で ErrorHandler 後にログ）。
//   - API 契約の型は api/openapi（oapi-codegen 生成）を利用する。
//
// 例:
//   - 認証: internal/delivery/restapi/v1/auth（Login / Register）が usecase/auth を呼び出す。
//
// ルール:
//   - ビジネスルール・永続化ロジックは持たない（変換と HTTP ステータス・エラー応答のマッピングのみ）。
//   - 命名は Request/Response（OpenAPI 型）を delivery で、Input/Output は usecase/query で使い分ける。
//   - ハンドラはエラー時 return err のみ。JSON 化は restapi/engine.go の HTTPErrorHandler が httperror.Encode で行う。
//   - ドメインエラー（domain/errors）等は httperror が HTTP ステータスと ErrorResponse にマッピングする。
//   - usecase / query の Execute には c.Request().Context() を渡す。
//   - ハンドラ（v1/auth 等）は domain や infra を直接参照しない（registry 経由の Interactor のみ呼ぶ）。
package delivery
