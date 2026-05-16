package config

import (
	"slices"

	"github.com/urfave/cli/v3"
)

// Server -.
type Server struct {
	APP  app
	HTTP http
	JWT  jwt
	DB   db
	Log  log
}

// Flags -.
func (c *Server) Flags() []cli.Flag {
	return slices.Concat(
		c.APP.flags(),
		c.HTTP.flags(),
		c.JWT.flags(),
		c.DB.flags(),
		c.Log.flags(),
	)
}

// Workflow -.
type Workflow struct {
	APP app
	Log log
}

// Flags -.
func (c *Workflow) Flags() []cli.Flag {
	return slices.Concat(
		c.APP.flags(),
		c.Log.flags(),
	)
}
