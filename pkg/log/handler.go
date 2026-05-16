package log

import (
	"context"
	"log/slog"
)

const labelUnknown = "unknown"

type wrapSlogHandler struct {
	slog.Handler
	contextKeys []string
}

func newWrapSlogHandler(handler slog.Handler, contextKeys ...string) slog.Handler {
	return &wrapSlogHandler{Handler: handler, contextKeys: contextKeys}
}

// override Handle method to add context keys
func (w *wrapSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, key := range w.contextKeys {
		value := ctx.Value(key)
		if value == nil {
			value = labelUnknown
		}
		record.AddAttrs(slog.Attr{Key: key, Value: slog.AnyValue(value)})
	}
	return w.Handler.Handle(ctx, record)
}
