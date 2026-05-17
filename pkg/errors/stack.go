package errors

import (
	"fmt"
	"runtime"
	"strings"
)

const modulePrefix = "go-clean-api/"

// StackFrames はログ出力向けにスタックを人が読める文字列の配列で返す。
func (e *Error[T]) StackFrames() []string {
	return formatStack(e.stack)
}

// StackTrace は StackFrames を改行区切りで連結した文字列を返す。
func (e *Error[T]) StackTrace() string {
	return strings.Join(e.StackFrames(), "\n")
}

func formatStack(pcs []uintptr) []string {
	if len(pcs) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs)
	out := make([]string, 0, len(pcs))

	for {
		frame, more := frames.Next()
		if shouldSkipStackFrame(frame.Function) {
			if !more {
				break
			}
			continue
		}

		out = append(out, formatStackFrame(frame))
		if !more {
			break
		}
	}
	return out
}

func shouldSkipStackFrame(function string) bool {
	// pkg/errors パッケージ内の New/Wrap/captureStack のみ除外（呼び出し元は残す）
	if strings.Contains(function, "/pkg/errors.captureStack") {
		return true
	}
	if !strings.Contains(function, "/pkg/errors.Type[") {
		return false
	}
	return strings.HasSuffix(function, ".New") || strings.HasSuffix(function, ".Wrap")
}

func formatStackFrame(frame runtime.Frame) string {
	file := frame.File
	if i := strings.Index(file, modulePrefix); i >= 0 {
		file = file[i:]
	}

	fn := frame.Function
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}

	return fmt.Sprintf("%s:%d %s", file, frame.Line, fn)
}
