package db

import "github.com/jackc/pgconn"

//nolint:gochecknoglobals
var (
	constraintEmail    = "users_email_key"
	constraintUsername = "users_username_key"
)

func pqErrConflictIn(err error, constraint string) bool {
	pqErr, ok := err.(*pgconn.PgError)
	return ok && pqErr.ConstraintName == constraint
}
