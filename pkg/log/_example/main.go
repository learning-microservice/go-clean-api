package main

import (
	"context"
	"os"

	"go-clean-api/pkg/log"
)

func main() {
	// standard log
	{
		logger, err := log.New(os.Stdout,
			log.WithLevel("info"),
			log.WithPretty(true),
		)
		if err != nil {
			panic(err)
		}

		logger.Info("hello world")
	}

	// with global attrs
	{
		logger, err := log.New(os.Stdout,
			log.WithLevel("info"),
			log.WithGlobalAttrs(
				"app", "go-clean-api",
				"version", "1.0.0",
				"environment", "dev",
			),
			log.WithPretty(true),
		)
		if err != nil {
			panic(err)
		}

		logger.Info("hello world")
	}

	// with context keys
	{
		ctxKeyRequestId := "request_id"
		ctxKeyUserId := "user_id"

		logger, err := log.New(os.Stdout,
			log.WithLevel("info"),
			log.WithContextKeys(ctxKeyRequestId, ctxKeyUserId),
			log.WithPretty(true),
		)
		if err != nil {
			panic(err)
		}

		ctx := context.WithValue(context.Background(), ctxKeyRequestId, "1234567890")
		ctx = context.WithValue(ctx, ctxKeyUserId, "U0000001")

		logger.InfoContext(ctx, "hello world")
	}
}
