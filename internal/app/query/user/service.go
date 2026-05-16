package user

import "go-clean-api/internal/app"

// UserQueryService はユーザ検索に関するインターフェースです。
// 複雑な検索条件や、複数のテーブルを結合した結果を返す場合に利用します。
type QueryService = app.Interactor[*SearchInput, []SearchOutput]
