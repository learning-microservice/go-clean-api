package config

import (
	"time"

	"github.com/urfave/cli/v3"
)

type db struct {
	Address         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (db *db) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "db.address",
			Value:       "root:password@tcp(localhost:3306)/localdb?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4",
			Usage:       "DB Address",
			Sources:     cli.EnvVars("DB_ADDRESS"),
			Destination: &db.Address,
		},
		&cli.IntFlag{
			Name:        "db.maxidleconns",
			Value:       3,
			Usage:       "DB Max Idle Conns",
			Sources:     cli.EnvVars("DB_MAX_IDLE_CONNS"),
			Destination: &db.MaxIdleConns,
		},
		&cli.IntFlag{
			Name:        "db.maxopenconns",
			Value:       3,
			Usage:       "DB Max Open Conns",
			Sources:     cli.EnvVars("DB_MAX_OPEN_CONNS"),
			Destination: &db.MaxOpenConns,
		},
		&cli.DurationFlag{
			Name:        "db.connmaxlifetime",
			Value:       5 * time.Minute,
			Usage:       "DB Conn Max Lifetime",
			Sources:     cli.EnvVars("DB_CONN_MAX_LIFETIME"),
			Destination: &db.ConnMaxLifetime,
		},
		&cli.DurationFlag{
			Name:        "db.connmaxidletime",
			Value:       1 * time.Minute,
			Usage:       "DB Conn Max Idle Time",
			Sources:     cli.EnvVars("DB_CONN_MAX_IDLE_TIME"),
			Destination: &db.ConnMaxIdleTime,
		},
	}
}
