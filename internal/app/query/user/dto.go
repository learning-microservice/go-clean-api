package user

import "time"

type SearchInput struct {
	NameQuery string
}

// SearchResponse は検索結果として返却するDTOです。
// UIの要件に合わせて、ドメインモデルとは独立して定義します。
type SearchOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"` // ドメインモデルにはないが、検索要件で必要になるかもしれない例
}
