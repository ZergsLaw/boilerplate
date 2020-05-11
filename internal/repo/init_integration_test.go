// +build integration

package repo_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/repo"
	"github.com/zergslaw/boilerplate/migration"
	"go.uber.org/zap"
)

var (
	Repo *repo.Repo

	timeoutConnect = time.Second * 5
)

func TestMain(m *testing.M) {
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

	dbConn, err := zergrepo.Connect(ctx, "postgres")
	if err != nil {
		log.Fatal(fmt.Errorf("connect UserRepo: %w", err))
	}

	metric := zergrepo.MustMetric("test", "repo")

	mapper := zergrepo.NewMapper(
		zergrepo.NewConvert(app.ErrNotFound, sql.ErrNoRows),
		zergrepo.PQConstraint(app.ErrEmailExist, repo.ConstraintEmail),
		zergrepo.PQConstraint(app.ErrUsernameExist, repo.ConstraintUsername),
	)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(fmt.Errorf("connect zap: %w", err))
	}

	Repo = repo.New(zergrepo.New(dbConn, logger, metric, mapper))

	os.Exit(m.Run())
}

func truncate() error {
	return Repo.Exec(ctx, "TRUNCATE users, sessions, notifications, recovery_code RESTART IDENTITY CASCADE")
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
