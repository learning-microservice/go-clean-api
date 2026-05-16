package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"

	"go-clean-api/config"
	"go-clean-api/internal/delivery/restapi"
	"go-clean-api/internal/registry"
	"go-clean-api/pkg/httpserver"
)

// NewCommand -.
func NewCommand() *cli.Command {
	// 1. 引数があればそれを使用
	// 2. 引数がなく環境変数があればそれを使用
	// 3. どちらもなければ DefaultValue を使用
	var cfg config.Server

	return &cli.Command{
		Name:  "server",
		Usage: "Command to start the API server",
		Flags: cfg.Flags(),
		Action: func(context.Context, *cli.Command) error {
			// setup registry
			reg, err := registry.New(&cfg)
			if err != nil {
				return fmt.Errorf("failed to setup registry: %w", err)
			}
			logger := reg.Logger

			// setup handler
			handler := restapi.New(&cfg, reg)

			// setup rest-api server
			server := httpserver.New(handler,
				httpserver.WithAddress(cfg.HTTP.Host, cfg.HTTP.Port),
				httpserver.WithReadHeaderTimeout(cfg.HTTP.ReadHeaderTimeout),
				httpserver.WithErrorLog(slog.NewLogLogger(logger.Handler(), slog.LevelError)),
				httpserver.OnBeforeStart(func(addr string) {
					logger.Info("starting http server...", "address", addr)
				}),
				httpserver.OnBeforeShutdown(func(addr string) {
					logger.Info("starting http server shutting down...", "address", addr)
				}),
				httpserver.OnAfterShutdown(func(addr string) {
					logger.Info("shutdown http server completed", "address", addr)
				}),
				httpserver.OnError(func(err error) {
					logger.Error("http server error", "error", err)
				}),
			)

			return server.ListenAndServe()
		},
	}
}
