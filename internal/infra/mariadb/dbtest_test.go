package mariadb_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go-clean-api/internal/domain/user"
	"go-clean-api/internal/infra/mariadb"
	"go-clean-api/pkg/sqldb"
)

const testDBName = "testdb"

var testDB *sqldb.Client

func schemaPath() string {
	return filepath.Join("..", "..", "..", "migrations", "schema.sql")
}

func runRepositoryTest(t *testing.T, fn func(ctx context.Context, repo user.Repository)) {
	t.Helper()

	if testing.Short() {
		t.Skip("integration test (skipped with -short)")
	}

	err := testDB.RunInRollbackTx(context.Background(), func(ctx context.Context) error {
		repo := mariadb.NewUserRepository(testDB)
		fn(ctx, repo)
		return nil
	})
	require.NoError(t, err)
}

// bcrypt ハッシュ形式のダミー（password_hash CHAR(60) 用）
func dummyPasswordHash() []byte {
	return []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
}
