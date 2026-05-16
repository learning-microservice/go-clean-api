package config

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v3"
)

type jwt struct {
	Secret      string
	TokenExpiry time.Duration
}

func (jwt *jwt) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "jwt.secret",
			Value:       "your-secret-key-change-in-production",
			Usage:       "JWT signing secret (development default; override in production)",
			Sources:     cli.EnvVars("JWT_SECRET"),
			Destination: &jwt.Secret,
		},
		&cli.DurationFlag{
			Name:        "jwt.token.expiry",
			Value:       24 * time.Hour,
			Usage:       "JWT Token Expiry",
			Sources:     cli.EnvVars("JWT_TOKEN_EXPIRY"),
			Destination: &jwt.TokenExpiry,
			Validator: func(d time.Duration) error {
				if d <= 0 {
					return fmt.Errorf("jwt token expiry must be greater than 0")
				}
				return nil
			},
		},
	}
}
