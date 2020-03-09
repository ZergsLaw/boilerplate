package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/flag"
	"github.com/zergslaw/boilerplate/internal/repo"
	"github.com/zergslaw/boilerplate/migration"
)

// Default value.
const (
	DBName = "postgres"
	DBUser = "postgres"
	DBPass = "postgres"
	DBHost = "localhost"
	DBPort = 5432

	DirMigrate = "migration"

	ConnectTimeout = time.Second * 5
)

// nolint:gochecknoglobals,gocritic
var (
	migrateDir = flag.NewStrFlag("dir-migrate", "goose migrations dir",
		flag.StrRequired(), flag.StrEnv("DIR_MIGRATE"), flag.StrDefault(DirMigrate))

	dbName = flag.NewStrFlag("db-name", "database name",
		flag.StrRequired(), flag.StrEnv("DB_NAME"), flag.StrAliases("N"), flag.StrDefault(DBName))
	dbUser = flag.NewStrFlag("db-user", "database user",
		flag.StrRequired(), flag.StrEnv("DB_USER"), flag.StrAliases("U"), flag.StrDefault(DBUser))
	dbPass = flag.NewStrFlag("db-pass", "database pass",
		flag.StrRequired(), flag.StrEnv("DB_PASS"), flag.StrAliases("P"), flag.StrDefault(DBPass))
	dbHost = flag.NewStrFlag("db-host", "database host",
		flag.StrRequired(), flag.StrEnv("DB_HOST"), flag.StrAliases("H"), flag.StrDefault(DBHost))
	dbPort = flag.NewIntFlag("db-port", "database port",
		flag.IntRequired(), flag.IntEnv("DB_PORT"), flag.IntAliases("p"), flag.IntDefault(DBPort))

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
		repo.Name(c.String(dbName.Name)),
		repo.User(c.String(dbUser.Name)),
		repo.Pass(c.String(dbPass.Name)),
		repo.Host(c.String(dbHost.Name)),
		repo.Port(c.Int(dbPort.Name)),
	)
}

func goose(ctx context.Context, dir, cmd string, opt ...repo.Option) error {
	return migration.Run(ctx, dir, cmd, opt...)
}
