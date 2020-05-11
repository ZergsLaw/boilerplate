// Package repo is an implements database interface.
package repo

import (
	"context"
	"database/sql"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	_ "github.com/lib/pq" // db driver.
)

// Repo is an implements app.UserRepo.
// Responsible for working with database.
type Repo struct {
	db *zergrepo.Repo
}

// Exec database query.
func (repo *Repo) Exec(ctx context.Context, query string) error {
	return repo.db.Do(func(db *sql.DB) error {
		_, err := db.ExecContext(ctx, query)
		return err
	})
}

// New creates and returns new app.UserRepo.
func New(repo *zergrepo.Repo) *Repo {
	return &Repo{db: repo}
}

// Constraint names.
const (
	ConstraintEmail    = "users_email_key"
	ConstraintUsername = "users_username_key"
)
