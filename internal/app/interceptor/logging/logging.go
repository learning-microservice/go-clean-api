package logging

import (
	"context"
	"log/slog"
	"time"

	"go-clean-api/internal/app"
)

func Wrap[I, O any](name string, logger *slog.Logger, nowFunc func() time.Time) app.Wrap[I, O] {
	return func(next app.Interactor[I, O]) app.Interactor[I, O] {
		return New(next,
			WithName[I, O](name),
			WithLogger[I, O](logger),
			WithNowFunc[I, O](nowFunc),
		)
	}
}

func New[I, O any](next app.Interactor[I, O], opts ...Option[I, O]) app.Interactor[I, O] {
	l := &logging[I, O]{next: next}

	// setup options
	for _, opt := range opts {
		opt(l)
	}

	if l.name == "" {
		l.name = "logging-interceptor"
	}

	if l.logger == nil {
		l.logger = slog.Default()
	}

	if l.nowFunc == nil {
		l.nowFunc = time.Now
	}

	return l
}

type logging[I, O any] struct {
	name    string
	logger  *slog.Logger
	nowFunc func() time.Time
	next    app.Interactor[I, O]
}

func (i *logging[I, O]) Execute(ctx context.Context, input I) (O, error) {
	start := i.nowFunc()
	output, err := i.next.Execute(ctx, input)

	var message, detail string
	if err == nil {
		message = "interactor executed successfully"
	} else {
		message = "error executing interactor"
		detail = err.Error()
	}

	// logging error
	// TODO: error type によって logging level を変えるが、現状は info でログ出力
	i.logger.InfoContext(ctx, message,
		"name", i.name,
		"duration", time.Since(start).String(),
		"detail", detail)

	return output, err
}
