package log

import (
	"errors"
	"io"
	"log/slog"
)

// New -.
func New(writer io.Writer, opts ...Option) (*slog.Logger, error) {
	cfg := config{
		option: slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		},
	}
	var errs []error
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	// setup pretty writer (for local development)
	if cfg.pretty {
		writer = &prettyJSONWriter{out: writer, indent: "  "}
	}

	// setup slog handler
	var handler slog.Handler = slog.NewJSONHandler(writer, &cfg.option)
	if len(cfg.contextKeys) > 0 {
		handler = newWrapSlogHandler(handler, cfg.contextKeys...)
	}

	logger := slog.New(handler)
	if len(cfg.attrs) > 0 {
		logger = logger.With(cfg.attrs...)
	}

	return logger, nil
}
