package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	"github.com/ZergsLaw/zerg-repo/zergrepo/cmd"
	"github.com/ZergsLaw/zerg-repo/zergrepo/core"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/log"
	"github.com/zergslaw/boilerplate/internal/repo"
	"github.com/zergslaw/boilerplate/internal/repo/migration"
)

// Default value.
const (
	ConnectTimeout = time.Second * 5
)

var (
	dbName = &cli.StringFlag{
		Name:       "db-name",
		Aliases:    []string{"n"},
		Usage:      "database name",
		EnvVars:    []string{"DB_NAME"},
		Value:      zergrepo.DBName,
		Required:   true,
		HasBeenSet: true,
	}
	operation = &cli.StringFlag{
		Name:    "operation",
		Aliases: []string{"o"},
		Usage: fmt.Sprintf("migration command one of (%s,%s,%s,%s,%s,%s)",
			core.Up, core.UpTo, core.UpOne, core.Down, core.DownTo, core.Reset),
		Required: true,
	}
	to = &cli.UintFlag{
		Name:    "to",
		Aliases: []string{"t"},
		Usage:   "on what element to migrate",
	}
	dbUser = &cli.StringFlag{
		Name:       "db-user",
		Aliases:    []string{"u"},
		Usage:      "database user",
		EnvVars:    []string{"DB_USER"},
		Value:      zergrepo.DBUser,
		Required:   true,
		HasBeenSet: true,
	}
	dbPass = &cli.StringFlag{
		Name:       "db-pass",
		Aliases:    []string{"p"},
		Usage:      "database password",
		EnvVars:    []string{"DB_PASS"},
		Value:      zergrepo.DBPassword,
		Required:   true,
		HasBeenSet: true,
	}
	dbHost = &cli.StringFlag{
		Name:       "db-host",
		Aliases:    []string{"H"},
		Usage:      "database host",
		EnvVars:    []string{"DB_HOST"},
		Value:      zergrepo.DBHost,
		Required:   true,
		HasBeenSet: true,
	}
	dbPort = &cli.IntFlag{
		Name:       "db-port",
		Aliases:    []string{"P"},
		Usage:      "database port",
		EnvVars:    []string{"DB_PORT"},
		Value:      zergrepo.DBPort,
		Required:   true,
		HasBeenSet: true,
	}

	Migrate = &cli.Command{
		Name:         "migrate",
		Aliases:      []string{"m"},
		Usage:        "causes the migration to the database.",
		UsageText:    "Migrate database schema.",
		BashComplete: cli.DefaultAppComplete,
		Action:       migrateAction,
		Flags:        []cli.Flag{dbName, dbUser, dbPass, dbHost, dbPort, operation},
	}
)

// Errors.
var (
	ErrUnknownCommand = errors.New("unknown cmd")
)

func migrateAction(c *cli.Context) error {
	ctxConnect, cancelConnect := context.WithTimeout(c.Context, ConnectTimeout)
	defer cancelConnect()

	dbConn, err := zergrepo.ConnectByCfg(ctxConnect, "postgres", zergrepo.Config{
		Host:     c.String(dbHost.Name),
		Port:     c.Int(dbPort.Name),
		User:     c.String(dbUser.Name),
		Password: c.String(dbPass.Name),
		DBName:   c.String(dbName.Name),
		SSLMode:  zergrepo.DBSSLMode,
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	zp := repo.Connect(dbConn, log.FromContext(c.Context).Named("zergrepo").Sugar(), c.App.Name)

	err = zergrepo.RegisterMetric(migration.Migrations...)
	if err != nil {
		return fmt.Errorf("register metric: %w", err)
	}

	command, err := parse(c.String(operation.Name))
	if err != nil {
		return fmt.Errorf("parse command: %w", err)
	}

	switch command {
	case core.Up:
		return zp.Up(c.Context)
	case core.UpOne:
		return zp.UpOne(c.Context)
	case core.UpTo:
		return zp.UpTo(c.Context, c.Uint(to.Name))
	case core.Down:
		return zp.Down(c.Context)
	case core.DownTo:
		return zp.DownTo(c.Context, c.Uint(to.Name))
	case core.Reset:
		return zp.Reset(c.Context)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownCommand, command)
	}
}

func parse(op string) (command core.MigrateCmd, err error) {
	switch op {
	case core.Up.String():
		command = core.Up
	case core.UpTo.String():
		command = core.UpTo
	case core.UpOne.String():
		command = core.UpOne
	case core.Down.String():
		command = core.Down
	case core.DownTo.String():
		command = core.DownTo
	case core.Reset.String():
		command = core.Reset
	default:
		return 0, cmd.ErrUnknownOperation
	}

	return command, nil
}
