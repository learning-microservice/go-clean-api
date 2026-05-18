package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqldblogger "github.com/simukti/sqldb-logger"
)

// DBとTxの両方が満たすインターフェース
type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Client struct {
	db *sql.DB
}

func NewClient(dsn string, opts ...Option) (*Client, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// setup default values
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// apply options
	for _, opt := range opts {
		opt(db)
	}

	// 4. sqldblogger で既存の *sql.DB をラップする
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	loggerAdapter := &slogAdapter{logger: jsonLogger}

	// OpenDriver を使ってロギング機能付きの *sql.DB を再生成
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)

	return &Client{db: db}, nil
}

func (c *Client) CurrentTx(ctx context.Context) DB {
	// Contextのトランザクションを取得した場合はそのまま返却
	if tx, ok := txFromContext(ctx); ok {
		return tx
	}
	return c.db
}

func (c *Client) IsDuplicateKeyError(err error) bool {
	return errorCode(err) == ErrCodeDuplicateKey
}

func (c *Client) IsForeignKeyError(err error) bool {
	return errorCode(err) == ErrCodeForeignKey
}

func (c *Client) Close() error {
	return c.db.Close()
}
