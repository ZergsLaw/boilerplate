package repo

import (
	"github.com/lib/pq"
)

const (
	constraintEmail    = "users_email_key"
	constraintUsername = "users_username_key"
)

func pqErrConflictIn(err error, constraint string) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Constraint == constraint
}
