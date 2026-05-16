//go:generate go tool mockgen -source=$GOFILE -package=$GOPACKAGE_test -destination=./mocks/$GOFILE
package auth

import (
	"go-clean-api/internal/app"
)

// --- Input Ports (Usecase Interface) ---

// LoginUsecase はユーザーログインを実行し、成功時にトークンを返します。
type LoginUsecase = app.Interactor[*LoginInput, *LoginOutput]

// RegisterUsecase はユーザー登録を実行し、成功時にユーザーIDを返します。
type RegisterUsecase = app.Interactor[*RegisterInput, *RegisterOutput]

// --- Output Ports (Repository Interface) ---

// TokenGenerator は指定されたユーザーIDに対するアクセストークンを生成します。
type TokenGenerator interface {
	Generate(userID string) (token string, err error)
}

// PasswordVerifier は平文パスワードと保存済みハッシュを照合します。
type PasswordVerifier interface {
	Verify(storedHash []byte, plain string) error
}

// PasswordHasher は平文パスワードをハッシュ化します。
type PasswordHasher interface {
	Hash(plain string) (hash []byte, err error)
}

// IDGenerator は新しいIDを生成します。
type IDGenerator interface {
	Generate() (id uint64, err error)
}
