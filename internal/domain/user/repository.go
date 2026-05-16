//go:generate go tool mockgen -source=$GOFILE -package=mocks -destination=./mocks/$GOFILE
package user

import (
	"context"
)

// ユーザリポジトリ
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	Save(ctx context.Context, user *User) (ID, error)
}
