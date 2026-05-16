package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"go-clean-api/cmd/app/server"
	"go-clean-api/cmd/app/workflow"
)

func main() {
	app := &cli.Command{
		Name: "go-clean-api",
		Commands: []*cli.Command{
			server.NewCommand(),
			workflow.NewCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
