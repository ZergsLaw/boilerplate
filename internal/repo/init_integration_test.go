// +build integration

package repo_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/repo"
	"github.com/zergslaw/boilerplate/migration"
)

var (
	Repo *repo.Repo

	timeoutConnect = time.Second * 5
)

func TestMain(m *testing.M) {
	repo.InitMetrics("test")

	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnect)
	defer cancel()

	resetDB := func() {
		err := migration.Run(ctx, "../../migration", "reset")
		if err != nil {
			log.Fatal(fmt.Errorf("migration: %w", err))
		}
	}
	// For convenient cleaning DB.
	resetDB()

	err := migration.Run(ctx, "../../migration", "up")
	if err != nil {
		log.Fatal(fmt.Errorf("migration: %w", err))
	}

	defer resetDB()

	dbConn, err := repo.Connect(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("connect UserRepo: %w", err))
	}

	Repo = repo.New(dbConn)

	os.Exit(m.Run())
}

func truncate() error {
	_, err := Repo.DB().Exec("TRUNCATE users, sessions, notifications, recovery_code RESTART IDENTITY CASCADE")
	return err
}

var (
	userGenerator = generatorUser()
	ctx           = context.Background()
	ip            = "192.100.10.4"
	origin        = app.Origin{
		IP:        net.ParseIP(ip),
		UserAgent: "UserAgent",
	}
)

func generatorUser() func() app.User {
	x := 0

	return func() app.User {
		x++
		return app.User{
			ID:        app.UserID(x),
			Email:     fmt.Sprintf("email%d@gmail.com", x),
			Username:  fmt.Sprintf("username%d", x),
			PassHash:  []byte("pass"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
}
