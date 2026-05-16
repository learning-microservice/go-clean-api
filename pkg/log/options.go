package log

import (
	"fmt"
	"log/slog"
	"strings"
)

// Option -.
type Option func(*config) error

type config struct {
	option      slog.HandlerOptions
	attrs       []any
	contextKeys []string
	pretty      bool
}

// WithLevel -.
func WithLevel(level string) Option {
	return func(c *config) error {
		switch strings.ToLower(level) {
		case "error":
			c.option.Level = slog.LevelError
		case "warn":
			c.option.Level = slog.LevelWarn
		case "info":
			c.option.Level = slog.LevelInfo
		case "debug":
			c.option.Level = slog.LevelDebug
		default:
			return fmt.Errorf("invalid log level: %s, use one of: error, warn, info, debug", level)
		}
		return nil
	}
}

// WithGlobalAttrs -.
func WithGlobalAttrs(attrs ...any) Option {
	return func(c *config) error {
		c.attrs = attrs
		return nil
	}
}

// WithContextKeys -.
func WithContextKeys(keys ...string) Option {
	return func(c *config) error {
		c.contextKeys = keys
		return nil
	}
}

// WithPretty -.
func WithPretty(pretty bool) Option {
	return func(c *config) error {
		c.pretty = pretty
		return nil
	}
}
