package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go-clean-api/internal/app/usecase/auth"
	mockauth "go-clean-api/internal/app/usecase/auth/mocks"
	domainerrors "go-clean-api/internal/domain/errors"
	"go-clean-api/internal/domain/user"
	usermocks "go-clean-api/internal/domain/user/mocks"
)

func TestRegisterUsecase_Execute(t *testing.T) {
	t.Parallel()

	const (
		name     = "テストユーザ"
		email    = "new@example.com"
		password = "plain-password"
	)

	input := &auth.RegisterInput{
		Name:          name,
		Email:         email,
		PlainPassword: password,
	}

	tests := []struct {
		name    string
		arrange func(
			repo *usermocks.MockRepository,
			hasher *mockauth.MockPasswordHasher,
		)
		assert func(t *testing.T, out *auth.RegisterOutput, err error)
	}{
		{
			name: "成功するとユーザーIDが返る",
			arrange: func(repo *usermocks.MockRepository, hasher *mockauth.MockPasswordHasher) {
				hasher.EXPECT().Hash(password).Return([]byte("hashed"), nil)
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).
					Return(user.NewID(42), nil)
			},
			assert: func(t *testing.T, out *auth.RegisterOutput, err error) {
				require.NoError(t, err)
				require.NotNil(t, out)
				assert.Equal(t, user.NewID(42), out.ID)
			},
		},
		{
			name: "既に登録済みの場合は重複エラー",
			arrange: func(repo *usermocks.MockRepository, hasher *mockauth.MockPasswordHasher) {
				hasher.EXPECT().Hash(password).Return([]byte("hashed"), nil)
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).
					Return(user.ID(0), domainerrors.TypeAlreadyExists.New("user already registered"))
			},
			assert: func(t *testing.T, out *auth.RegisterOutput, err error) {
				require.Error(t, err)
				assert.Nil(t, out)
				assert.True(t, domainerrors.TypeAlreadyExists.Is(err))
			},
		},
		{
			name: "ハッシュ化失敗は予期しないエラー",
			arrange: func(_ *usermocks.MockRepository, hasher *mockauth.MockPasswordHasher) {
				hasher.EXPECT().Hash(password).Return(nil, errors.New("hash error"))
			},
			assert: func(t *testing.T, out *auth.RegisterOutput, err error) {
				require.Error(t, err)
				assert.Nil(t, out)
				assert.True(t, domainerrors.TypeUnexpected.Is(err))
			},
		},
		{
			name: "保存失敗は予期しないエラー",
			arrange: func(repo *usermocks.MockRepository, hasher *mockauth.MockPasswordHasher) {
				hasher.EXPECT().Hash(password).Return([]byte("hashed"), nil)
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).
					Return(user.ID(0), errors.New("insert failed"))
			},
			assert: func(t *testing.T, out *auth.RegisterOutput, err error) {
				require.Error(t, err)
				assert.Nil(t, out)
				assert.True(t, domainerrors.TypeUnexpected.Is(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			repo := usermocks.NewMockRepository(ctrl)
			hasher := mockauth.NewMockPasswordHasher(ctrl)

			tt.arrange(repo, hasher)

			uc := auth.NewRegisterUsecase(repo, hasher)
			out, err := uc.Execute(context.Background(), input)

			tt.assert(t, out, err)
		})
	}
}
