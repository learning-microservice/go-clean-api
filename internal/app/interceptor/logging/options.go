package logging

import (
	"log/slog"
	"time"
)

type Option[I, O any] func(*logging[I, O])

func WithName[I, O any](name string) Option[I, O] {
	return func(l *logging[I, O]) {
		l.name = name
	}
}

func WithLogger[I, O any](logger *slog.Logger) Option[I, O] {
	return func(l *logging[I, O]) {
		l.logger = logger
	}
}

func WithNowFunc[I, O any](nowFunc func() time.Time) Option[I, O] {
	return func(l *logging[I, O]) {
		l.nowFunc = nowFunc
	}
}
