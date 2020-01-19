// Package migration allows for migration and rejection of tables in the database.
package migration

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"github.com/zergslaw/users/internal/db"
)

// nolint:gochecknoglobals
var (
	gooseMu sync.Mutex
	log     logrus.FieldLogger = logrus.New().WithField("package", "db")
)

func warnIfFail(fn func() error) {
	if err := fn(); err != nil {
		log.Warn(err)
	}
}

// Run executes goose command. It also enforce "fix" after "create".
func Run(ctx context.Context, dir string, command string, options ...db.Option) error {
	gooseMu.Lock()
	defer gooseMu.Unlock()

	dbConn, err := db.Connect(ctx, options...)
	if err != nil {
		return err
	}
	defer warnIfFail(dbConn.Close)

	cmdArgs := strings.Fields(command)
	cmd, args := cmdArgs[0], cmdArgs[1:]
	err = goose.Run(cmd, dbConn.DB, dir, args...)
	if err == nil && cmd == "create" {
		err = goose.Run("fix", dbConn.DB, dir)
	}
	if err != nil {
		return fmt.Errorf("goose.Run %q: %w", command, err)
	}
	return nil
}
