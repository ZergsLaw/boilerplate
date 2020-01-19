// Package db is an implements app.Repo.
package db

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // driver fot postgres.
	"github.com/jmoiron/sqlx"
	"github.com/zergslaw/users/internal/app"
)

// Repo is an implements app.Repo.
// Responsible for working with database.
type Repo struct {
	db *sqlx.DB
}

// Close closes database connections.
func (repo *Repo) Close() error {
	return repo.db.Close()
}

// New creates and returns new app.Repo.
func New(conn *sqlx.DB) app.Repo {
	return &Repo{db: conn}
}

const (
	dbMaxOpenConns  = 30 // about ⅓ of server's max_connections
	dbParallelConns = 5  // a bit more than average
)

// Connect to database by options.
func Connect(ctx context.Context, options ...Option) (*sqlx.DB, error) {
	opt := defaultConfig()

	for i := range options {
		options[i](opt)
	}

	dbConn, err := sqlx.Open("pgx", opt.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx open: %w", err)
	}
	dbConn.SetMaxOpenConns(dbMaxOpenConns)
	dbConn.SetMaxIdleConns(dbParallelConns)

	err = dbConn.PingContext(ctx)
	for err != nil {
		nextErr := dbConn.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			if errClose := dbConn.Close(); errClose != nil {
				return nil, fmt.Errorf("db ping: %w, db close: %s", err, errClose)
			}
			return nil, fmt.Errorf("db ping: %w", err)
		}
		err = nextErr
	}

	return dbConn, nil
}

func (repo *Repo) execFunc(f func(db *sqlx.DB) error) (err error) {
	methodName, methodDone := methodMetrics(1)
	defer methodDone(&err)

	err = f(repo.db)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	return nil
}
