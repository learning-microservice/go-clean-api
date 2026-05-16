package config

import (
	"github.com/urfave/cli/v3"
)

type app struct {
	Name    string
	Version string
}

func (a *app) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "app.name",
			Value:       "go-clean-api",
			Usage:       "Application Name",
			Sources:     cli.EnvVars("APP_NAME"),
			Destination: &a.Name,
		},
		&cli.StringFlag{
			Name:        "app.version",
			Value:       "unknown",
			Usage:       "Application Version",
			Sources:     cli.EnvVars("APP_VERSION"),
			Destination: &a.Version,
		},
	}
}
