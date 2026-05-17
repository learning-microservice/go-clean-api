package errors_test

import (
	"strings"
	"testing"

	pkgErrors "go-clean-api/pkg/errors"
)

func TestStackFrames_readableAndSkipsPkgErrors(t *testing.T) {
	t.Parallel()

	typ := pkgErrors.Type[string](400, "test")
	err := typ.New("boom")

	frames := err.StackFrames()
	if len(frames) == 0 {
		t.Fatal("expected at least one stack frame")
	}

	for _, frame := range frames {
		if strings.Contains(frame, "/pkg/errors.Type[") || strings.Contains(frame, "captureStack") {
			t.Errorf("internal pkg/errors machinery frame should be skipped: %q", frame)
		}
		if !strings.Contains(frame, ":") {
			t.Errorf("frame should contain file:line: %q", frame)
		}
	}

	trace := err.StackTrace()
	if trace == "" {
		t.Fatal("expected non-empty stack trace string")
	}
	if len(frames) > 1 && !strings.Contains(trace, "\n") {
		t.Errorf("expected newline-separated trace, got %q", trace)
	}
}

func TestStackFrames_wrapUsesOuterStack(t *testing.T) {
	t.Parallel()

	typ := pkgErrors.Type[string](400, "test")
	inner := typ.New("inner")
	outer := typ.Wrap(inner, "outer")

	if len(outer.StackFrames()) == 0 {
		t.Fatal("expected stack on wrapped error")
	}
}
