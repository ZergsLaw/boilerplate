// +build integration

package db_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/db"
	"github.com/zergslaw/users/migration"
)

var (
	Repo app.Repo

	timeoutConnect = time.Second * 5
)

func TestMain(m *testing.M) {
	db.InitMetrics("test")

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

	dbConn, err := db.Connect(ctx)
	if err != nil {
		panic(fmt.Errorf("connect Repo: %w", err))
	}

	Repo = db.New(dbConn)

	os.Exit(m.Run())
}
