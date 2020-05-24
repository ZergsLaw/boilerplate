// Package repo is an implements database interface.
package repo

import (
	"database/sql"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // db driver.
	"github.com/zergslaw/boilerplate/internal/app"
)

// Connect create new instance *zergrepo.Repo.
func Connect(db *sqlx.DB, logger zergrepo.Logger, namespace string) *zergrepo.Repo {
	// Constraint names.
	const (
		ConstraintEmail    = "users_email_key"
		ConstraintUsername = "users_username_key"
	)

	metric := zergrepo.MustMetric(namespace, "repo")
	mapper := zergrepo.NewMapper(
		zergrepo.NewConvert(app.ErrNotFound, sql.ErrNoRows),
		zergrepo.PQConstraint(app.ErrEmailExist, ConstraintEmail),
		zergrepo.PQConstraint(app.ErrUsernameExist, ConstraintUsername),
	)

	return zergrepo.New(db, logger, metric, mapper)
}

var _ app.SessionRepo = &Repo{}
var _ app.UserRepo = &Repo{}
var _ app.WAL = &Repo{}
var _ app.CodeRepo = &Repo{}

// Repo is an implements app.UserRepo.
// Responsible for working with database.
type Repo struct {
	db *zergrepo.Repo
}

// New creates and returns new app.UserRepo.
func New(repo *zergrepo.Repo) *Repo {
	return &Repo{db: repo}
}
