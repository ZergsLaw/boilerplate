// Package migration allows for migration and rejection of tables in the database.
package migration

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pressly/goose"
	"github.com/zergslaw/boilerplate/internal/repo"
)

// nolint:gochecknoglobals,gocritic
var (
	gooseMu sync.Mutex
)

// Run executes goose command. It also enforce "fix" after "create".
func Run(ctx context.Context, dir string, command string, options ...repo.Option) error {
	gooseMu.Lock()
	defer gooseMu.Unlock()

	dbConn, err := repo.Connect(ctx, options...)
	if err != nil {
		return err
	}

	cmdArgs := strings.Fields(command)
	cmd, args := cmdArgs[0], cmdArgs[1:]
	err = goose.Run(cmd, dbConn, dir, args...)
	if err == nil && cmd == "create" {
		err = goose.Run("fix", dbConn, dir)
	}
	if err != nil {
		return fmt.Errorf("goose.Run %q: %w", command, err)
	}

	return dbConn.Close()
}
