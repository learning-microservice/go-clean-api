package config

import (
	"github.com/urfave/cli/v3"
)

type log struct {
	Level  string
	Pretty bool
}

func (log *log) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log.level",
			Value:       "DEBUG",
			Usage:       "Logging Level",
			Sources:     cli.EnvVars("LOG_LEVEL"),
			Destination: &log.Level,
		},
		&cli.BoolFlag{
			Name:        "log.pretty",
			Value:       false,
			Usage:       "Logging Pretty JSON",
			Sources:     cli.EnvVars("LOG_PRETTY"),
			Destination: &log.Pretty,
		},
	}
}
