package mariadb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "go-clean-api/internal/domain/errors"
	"go-clean-api/internal/domain/user"
)

func TestUserRepository_FindByEmail(t *testing.T) {
	const email = "find@example.com"

	tests := []struct {
		name   string
		seed   func(ctx context.Context, repo user.Repository)
		email  string
		assert func(t *testing.T, got *user.User, err error)
	}{
		{
			name:  "未登録のメールアドレスは NotFound",
			email: "none@example.com",
			assert: func(t *testing.T, got *user.User, err error) {
				require.Error(t, err)
				assert.Nil(t, got)
				assert.True(t, domainerrors.TypeNotFound.Is(err))
			},
		},
		{
			name: "登録済みユーザをメールで取得できる",
			seed: func(ctx context.Context, repo user.Repository) {
				u := user.New("太郎", email, dummyPasswordHash())
				_, err := repo.Save(ctx, u)
				require.NoError(t, err)
			},
			email: email,
			assert: func(t *testing.T, got *user.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, "太郎", got.Name())
				assert.Equal(t, email, got.Email())
				assert.Equal(t, dummyPasswordHash(), got.PasswordHash())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runRepositoryTest(t, func(ctx context.Context, repo user.Repository) {
				if tt.seed != nil {
					tt.seed(ctx, repo)
				}
				got, err := repo.FindByEmail(ctx, tt.email)
				tt.assert(t, got, err)
			})
		})
	}
}

func TestUserRepository_Save(t *testing.T) {
	tests := []struct {
		name   string
		run    func(ctx context.Context, repo user.Repository)
		assert func(t *testing.T, repo user.Repository, ctx context.Context)
	}{
		{
			name: "新規ユーザを INSERT できる",
			run: func(ctx context.Context, repo user.Repository) {
				u := user.New("新規", "new@example.com", dummyPasswordHash())
				id, err := repo.Save(ctx, u)
				require.NoError(t, err)
				assert.NotZero(t, id.Int())

				got, err := repo.FindByEmail(ctx, "new@example.com")
				require.NoError(t, err)
				assert.Equal(t, "新規", got.Name())
				assert.Equal(t, id, got.ID())
			},
		},
		{
			name: "既存ユーザを UPDATE できる",
			run: func(ctx context.Context, repo user.Repository) {
				u := user.New("更新前", "update@example.com", dummyPasswordHash())
				id, err := repo.Save(ctx, u)
				require.NoError(t, err)

				updated := user.Reconstruct(id, "更新後", "update@example.com", dummyPasswordHash())
				_, err = repo.Save(ctx, updated)
				require.NoError(t, err)
			},
			assert: func(t *testing.T, repo user.Repository, ctx context.Context) {
				got, err := repo.FindByEmail(ctx, "update@example.com")
				require.NoError(t, err)
				assert.Equal(t, "更新後", got.Name())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runRepositoryTest(t, func(ctx context.Context, repo user.Repository) {
				tt.run(ctx, repo)
				if tt.assert != nil {
					tt.assert(t, repo, ctx)
				}
			})
		})
	}
}
