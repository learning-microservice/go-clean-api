package sqldb

import (
	"database/sql"
	"time"
)

type Option func(*sql.DB)

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(maxIdleConns)
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(maxOpenConns)
	}
}

func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(connMaxLifetime)
	}
}

func WithConnMaxIdleTime(connMaxIdleTime time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxIdleTime(connMaxIdleTime)
	}
}
