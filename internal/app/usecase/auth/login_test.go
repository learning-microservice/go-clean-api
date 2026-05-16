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

func TestLoginUsecase_Execute(t *testing.T) {
	t.Parallel()

	const (
		email    = "user@example.com"
		password = "plain-password"
	)

	existingUser := user.Reconstruct(
		user.NewID(1),
		"テストユーザ",
		email,
		[]byte("stored-hash"),
	)

	tests := []struct {
		name    string
		input   *auth.LoginInput
		arrange func(
			repo *usermocks.MockRepository,
			verifier *mockauth.MockPasswordVerifier,
			generator *mockauth.MockTokenGenerator,
		)
		assert func(t *testing.T, out *auth.LoginOutput, err error)
	}{
		{
			name:  "成功するとトークンが返る",
			input: &auth.LoginInput{Email: email, PlainPassword: password},
			arrange: func(
				repo *usermocks.MockRepository,
				verifier *mockauth.MockPasswordVerifier,
				generator *mockauth.MockTokenGenerator,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), email).Return(existingUser, nil)
				verifier.EXPECT().Verify([]byte("stored-hash"), password).Return(nil)
				generator.EXPECT().Generate("1").Return("jwt-token", nil)
			},
			assert: func(t *testing.T, out *auth.LoginOutput, err error) {
				require.NoError(t, err)
				require.NotNil(t, out)
				assert.Equal(t, "jwt-token", out.Token)
			},
		},
		{
			name:  "ユーザが見つからない場合は認証エラー",
			input: &auth.LoginInput{Email: email, PlainPassword: password},
			arrange: func(
				repo *usermocks.MockRepository,
				_ *mockauth.MockPasswordVerifier,
				_ *mockauth.MockTokenGenerator,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), email).
					Return(nil, domainerrors.TypeNotFound.New("user not found"))
			},
			assert: func(t *testing.T, out *auth.LoginOutput, err error) {
				require.Error(t, err)
				assert.Nil(t, out)
				assert.True(t, domainerrors.TypeInvalidCredentials.Is(err))
			},
		},
		{
			name:  "パスワード不一致は認証エラー",
			input: &auth.LoginInput{Email: email, PlainPassword: password},
			arrange: func(
				repo *usermocks.MockRepository,
				verifier *mockauth.MockPasswordVerifier,
				_ *mockauth.MockTokenGenerator,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), email).Return(existingUser, nil)
				verifier.EXPECT().Verify([]byte("stored-hash"), password).
					Return(errors.New("password mismatch"))
			},
			assert: func(t *testing.T, out *auth.LoginOutput, err error) {
				require.Error(t, err)
				assert.Nil(t, out)
				assert.True(t, domainerrors.TypeInvalidCredentials.Is(err))
			},
		},
		{
			name:  "トークン生成失敗は予期しないエラー",
			input: &auth.LoginInput{Email: email, PlainPassword: password},
			arrange: func(
				repo *usermocks.MockRepository,
				verifier *mockauth.MockPasswordVerifier,
				generator *mockauth.MockTokenGenerator,
			) {
				repo.EXPECT().FindByEmail(gomock.Any(), email).Return(existingUser, nil)
				verifier.EXPECT().Verify([]byte("stored-hash"), password).Return(nil)
				generator.EXPECT().Generate("1").Return("", errors.New("jwt error"))
			},
			assert: func(t *testing.T, out *auth.LoginOutput, err error) {
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
			verifier := mockauth.NewMockPasswordVerifier(ctrl)
			generator := mockauth.NewMockTokenGenerator(ctrl)

			tt.arrange(repo, verifier, generator)

			uc := auth.NewLoginUsecase(repo, generator, verifier)
			out, err := uc.Execute(context.Background(), tt.input)

			tt.assert(t, out, err)
		})
	}
}
