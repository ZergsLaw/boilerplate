// +build integration

package repo_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zergslaw/boilerplate/internal/repo"
	"github.com/zergslaw/boilerplate/migration"
)

var (
	Repo *repo.Repo

	timeoutConnect = time.Second * 1000
)

func TestMain(m *testing.M) {
	repo.InitMetrics("test")

	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnect)
	defer cancel()

	resetDB := func() {
		err := migration.Run(ctx, "../../migration", "reset")
		if err != nil {
			panic(fmt.Errorf("migration: %w", err))
		}
	}
	// For convenient cleaning DB.
	resetDB()

	err := migration.Run(ctx, "../../migration", "up")
	if err != nil {
		panic(fmt.Errorf("migration: %w", err))
	}

	defer resetDB()

	dbConn, err := repo.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("connect UserRepo: %w", err))
	}

	Repo = repo.New(dbConn)

	os.Exit(m.Run())
}
