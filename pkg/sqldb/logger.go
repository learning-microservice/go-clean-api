package sqldb

import (
	"context"
	"log/slog"

	sqldblogger "github.com/simukti/sqldb-logger"
)

type slogAdapter struct {
	logger *slog.Logger
}

func (s *slogAdapter) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	// ログデータ（map）を slog.Attr のスライスに変換
	attrs := make([]any, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	// sqldblogger のログレベルを slog のレベルにマッピングして出力
	switch level {
	case sqldblogger.LevelError:
		s.logger.ErrorContext(ctx, msg, attrs...)
	case sqldblogger.LevelInfo:
		s.logger.InfoContext(ctx, msg, attrs...)
	case sqldblogger.LevelDebug:
		s.logger.DebugContext(ctx, msg, attrs...)
	case sqldblogger.LevelTrace:
		// Traceはslog.LevelDebugで代用
		s.logger.DebugContext(ctx, msg, attrs...)
	}
}
