package main

import (
	"context"
	"errors"
	"strings"
	"time"

	zergrepo "github.com/ZergsLaw/zerg-repo"

	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/flag"
	"github.com/zergslaw/boilerplate/migration"
)

// Default value.
const (
	DirMigrate = "migration"

	ConnectTimeout = time.Second * 5
)

var (
	migrateDir = flag.NewStrFlag("dir-migrate", "goose migrations dir",
		flag.StrRequired(), flag.StrEnv("DIR_MIGRATE"), flag.StrDefault(DirMigrate))

	dbName = flag.NewStrFlag("db-name", "database name",
		flag.StrRequired(), flag.StrEnv("DB_NAME"), flag.StrAliases("N"), flag.StrDefault(zergrepo.DBName))
	dbUser = flag.NewStrFlag("db-user", "database user",
		flag.StrRequired(), flag.StrEnv("DB_USER"), flag.StrAliases("U"), flag.StrDefault(zergrepo.DBUser))
	dbPass = flag.NewStrFlag("db-pass", "database pass",
		flag.StrRequired(), flag.StrEnv("DB_PASS"), flag.StrAliases("P"), flag.StrDefault(zergrepo.DBPassword))
	dbHost = flag.NewStrFlag("db-host", "database host",
		flag.StrRequired(), flag.StrEnv("DB_HOST"), flag.StrAliases("H"), flag.StrDefault(zergrepo.DBHost))
	dbPort = flag.NewIntFlag("db-port", "database port",
		flag.IntRequired(), flag.IntEnv("DB_PORT"), flag.IntAliases("p"), flag.IntDefault(zergrepo.DBPort))

	migrate = &cli.Command{
		Name:         "migrate",
		Aliases:      []string{"m"},
		Usage:        "causes the migration to the database.",
		UsageText:    "Migrate database schema.",
		BashComplete: cli.DefaultAppComplete,
		Action:       migrateAction,
		Flags:        []cli.Flag{dbName, dbUser, dbPass, dbHost, dbPort, migrateDir},
	}
)

// Errors.
var (
	errGooseCommandRequired = errors.New("goose command is required")
)

func migrateAction(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return errGooseCommandRequired
	}

	command := strings.Join(c.Args().Slice(), " ")

	return goose(c.Context, c.String(migrateDir.Name), command,
		zergrepo.Name(c.String(dbName.Name)),
		zergrepo.User(c.String(dbUser.Name)),
		zergrepo.Pass(c.String(dbPass.Name)),
		zergrepo.Host(c.String(dbHost.Name)),
		zergrepo.Port(c.Int(dbPort.Name)),
	)
}

func goose(ctx context.Context, dir, cmd string, opt ...zergrepo.Option) error {
	return migration.Run(ctx, dir, cmd, opt...)
}
