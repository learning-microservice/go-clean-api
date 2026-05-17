package mariadb_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"

	"go-clean-api/pkg/sqldb"
)

const (
	testContainerName = "go-clean-api-mariadb-itest"
	testMariaDBImage  = "mariadb:11"
	testRootPassword  = "password"

	envTestDBAddress = "TEST_DB_ADDRESS"
	envReuseDisable  = "TESTCONTAINERS_REUSE_DISABLE"
)

var (
	mariadbContainer testcontainers.Container
	managedContainer bool
)

func TestMain(m *testing.M) {
	if isGoTestShort() {
		os.Exit(m.Run())
	}

	ctx := context.Background()

	if dsn := strings.TrimSpace(os.Getenv(envTestDBAddress)); dsn != "" {
		if err := setupExternalTestDB(ctx, dsn); err != nil {
			log.Printf("failed to setup external test database: %v", err)
			os.Exit(1)
		}
	} else {
		if err := setupTestcontainersDB(ctx); err != nil {
			log.Printf("failed to setup testcontainers database: %v", err)
			os.Exit(1)
		}
		managedContainer = true
	}

	code := m.Run()

	if err := testDB.Close(); err != nil {
		log.Printf("failed to close sqldb client: %v", err)
	}

	if managedContainer && mariadbContainer != nil && !containerReuseEnabled() {
		if err := mariadbContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate mariadb container: %v", err)
		}
	}

	os.Exit(code)
}

func setupExternalTestDB(ctx context.Context, dsn string) error {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("parse TEST_DB_ADDRESS: %w", err)
	}
	if cfg.DBName == "" {
		return fmt.Errorf("TEST_DB_ADDRESS must include a database name")
	}
	if err := validateTestDatabaseName(cfg.DBName); err != nil {
		return err
	}

	adminCfg := *cfg
	adminCfg.DBName = ""
	adminCfg.MultiStatements = true
	adminDSN := adminCfg.FormatDSN()

	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		return fmt.Errorf("open admin database: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping admin database: %w", err)
	}

	if err := resetDatabase(ctx, adminDB, cfg.DBName); err != nil {
		return err
	}

	client, err := sqldb.NewClient(dsn)
	if err != nil {
		return fmt.Errorf("create sqldb client: %w", err)
	}
	testDB = client

	return nil
}

func setupTestcontainersDB(ctx context.Context) error {
	configureTestcontainersReuse()

	container, host, port, err := startMariaDB(ctx)
	if err != nil {
		return err
	}
	mariadbContainer = container

	adminDB, err := sql.Open("mysql", rootDSN(host, port))
	if err != nil {
		return fmt.Errorf("open admin database: %w", err)
	}
	defer adminDB.Close()

	if err := adminDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping admin database: %w", err)
	}

	if err := resetDatabase(ctx, adminDB, testDBName); err != nil {
		return err
	}

	client, err := sqldb.NewClient(testDBDSN(host, port))
	if err != nil {
		return fmt.Errorf("create sqldb client: %w", err)
	}
	testDB = client

	return nil
}

func validateTestDatabaseName(dbName string) error {
	if dbName != testDBName {
		return fmt.Errorf(
			"integration tests may only reset database %q (got %q); use a dedicated test database",
			testDBName,
			dbName,
		)
	}
	return nil
}

func startMariaDB(ctx context.Context) (testcontainers.Container, string, string, error) {
	container, err := mariadb.Run(ctx,
		testMariaDBImage,
		mariadb.WithUsername("root"),
		mariadb.WithPassword(testRootPassword),
		mariadb.WithDatabase(testDBName),
		testcontainers.WithReuseByName(testContainerName),
	)
	if err != nil {
		return nil, "", "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", err
	}

	mappedPort, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, "", "", err
	}

	return container, host, mappedPort.Port(), nil
}

func resetDatabase(ctx context.Context, admin *sql.DB, dbName string) error {
	if _, err := admin.ExecContext(ctx, "DROP DATABASE IF EXISTS `"+dbName+"`"); err != nil {
		return fmt.Errorf("drop database: %w", err)
	}

	if _, err := admin.ExecContext(ctx,
		"CREATE DATABASE `"+dbName+"` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
	); err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	if _, err := admin.ExecContext(ctx, "USE `"+dbName+"`"); err != nil {
		return fmt.Errorf("use database: %w", err)
	}

	ddl, err := os.ReadFile(schemaPath())
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	if _, err := admin.ExecContext(ctx, string(ddl)); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
}

func rootDSN(host, port string) string {
	return fmt.Sprintf(
		"root:%s@tcp(%s)/?parseTime=true&multiStatements=true&loc=Asia%%2FTokyo&charset=utf8mb4",
		testRootPassword,
		net.JoinHostPort(host, port),
	)
}

func testDBDSN(host, port string) string {
	return fmt.Sprintf(
		"root:%s@tcp(%s)/%s?parseTime=true&loc=Asia%%2FTokyo&charset=utf8mb4",
		testRootPassword,
		net.JoinHostPort(host, port),
		testDBName,
	)
}

// configureTestcontainersReuse は reuse をデフォルト ON にする。
// testcontainers-go が参照する TESTCONTAINERS_REUSE_ENABLE へ変換する。
func configureTestcontainersReuse() {
	if isReuseDisabled() {
		_ = os.Setenv("TESTCONTAINERS_REUSE_ENABLE", "false")
		return
	}
	_ = os.Setenv("TESTCONTAINERS_REUSE_ENABLE", "true")
}

func containerReuseEnabled() bool {
	return !isReuseDisabled()
}

func isReuseDisabled() bool {
	if v := os.Getenv(envReuseDisable); v == "true" || v == "1" {
		return true
	}
	// testcontainers 標準変数で明示オフされた場合も尊重
	if v := os.Getenv("TESTCONTAINERS_REUSE_ENABLE"); v == "false" || v == "0" {
		return true
	}
	return false
}

// TestMain 実行時点では testing.Short() が未初期化のため、go test のフラグを直接見る。
func isGoTestShort() bool {
	for _, arg := range os.Args {
		if arg == "-test.short" || strings.HasPrefix(arg, "-test.short=") {
			return true
		}
	}
	return false
}
