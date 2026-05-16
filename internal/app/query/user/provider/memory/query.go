package memory

import (
	"context"

	"go-clean-api/internal/app/query/user"
)

type userQuery struct {
	// 本来はメモリ上のデータソースを参照しますが、
	// ここではプロトタイプとして固定値を返却する例とします。
}

func NewUserQuery() user.QueryService {
	return &userQuery{}
}

func (q *userQuery) Execute(_ context.Context, _ *user.SearchInput) ([]user.SearchOutput, error) {
	// モックデータの返却
	return []user.SearchOutput{
		{
			ID:    "1",
			Name:  "Test User",
			Email: "test@example.com",
		},
	}, nil
}
