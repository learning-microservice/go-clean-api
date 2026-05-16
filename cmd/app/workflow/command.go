package workflow

import (
	"context"

	"github.com/urfave/cli/v3"

	"go-clean-api/config"
)

// NewCommand -.
func NewCommand() *cli.Command {
	var cfg config.Workflow

	return &cli.Command{
		Name:  "workflow",
		Usage: "Command to start the workflow",
		Flags: cfg.Flags(),
		Action: func(context.Context, *cli.Command) error {
			// TODO: Implement the workflow
			return nil
		},
	}
}
