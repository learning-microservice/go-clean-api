package config

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v3"
)

// http -.
type http struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	GracefulTimeout   time.Duration
}

func (http *http) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http.host",
			Value:       "",
			Usage:       "HTTP Host",
			Sources:     cli.EnvVars("HTTP_HOST"),
			Destination: &http.Host,
		},
		&cli.IntFlag{
			Name:        "http.port",
			Value:       8080,
			Usage:       "HTTP Port",
			Sources:     cli.EnvVars("HTTP_PORT"),
			Destination: &http.Port,
			Validator: func(val int) error {
				if val < 1024 || val > 65535 {
					return fmt.Errorf("http port must be in range 1024-65535")
				}
				return nil
			},
		},
		&cli.DurationFlag{
			Name:        "http.readheadertimeout",
			Value:       5 * time.Second,
			Usage:       "HTTP ReadHeaderTimeout",
			Sources:     cli.EnvVars("HTTP_READ_HEADER_TIMEOUT"),
			Destination: &http.ReadHeaderTimeout,
		},
		&cli.DurationFlag{
			Name:        "http.readtimeout",
			Value:       5 * time.Second,
			Usage:       "HTTP ReadTimeout",
			Sources:     cli.EnvVars("HTTP_READ_TIMEOUT"),
			Destination: &http.ReadTimeout,
		},
		&cli.DurationFlag{
			Name:        "http.writetimeout",
			Value:       10 * time.Second,
			Usage:       "HTTP WriteTimeout",
			Sources:     cli.EnvVars("HTTP_WRITE_TIMEOUT"),
			Destination: &http.WriteTimeout,
		},
		&cli.DurationFlag{
			Name:        "http.idletimeout",
			Value:       60 * time.Second,
			Usage:       "HTTP IdleTimeout",
			Sources:     cli.EnvVars("HTTP_IDLE_TIMEOUT"),
			Destination: &http.IdleTimeout,
		},
		&cli.DurationFlag{
			Name:        "http.gracefultimeout",
			Value:       60 * time.Second,
			Usage:       "HTTP GracefulTimeout",
			Sources:     cli.EnvVars("HTTP_GRACEFUL_TIMEOUT"),
			Destination: &http.GracefulTimeout,
		},
	}
}
