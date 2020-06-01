package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	migrate "github.com/ZergsLaw/zerg-repo/zergrepo/cmd"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/cmd"
	"github.com/zergslaw/boilerplate/internal/api/web"
	"github.com/zergslaw/boilerplate/internal/api/web/generated/restapi"
	"github.com/zergslaw/boilerplate/internal/log"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	ver    string
	app    = &cli.App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "Boilerplate application.",
		BashComplete: cli.DefaultAppComplete,
		Writer:       os.Stdout,
		Commands:     []*cli.Command{cmd.Version, migrate.Migrate, cmd.Serve},
	}
)

func main() {
	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(app.Name, "_")
	web.InitMetrics(namespace, restapi.FlatSwaggerJSON)

	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init logger: %w", err))
	}

	ctx := log.SetContext(context.Background(), logger)
	ctx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() { <-signals; cancel() }()
	go forceShutdown(ctx)

	err = app.RunContext(ctx, os.Args)
	if err != nil {
		logger.Fatal("run service", zap.Error(err))
	}
}

func forceShutdown(ctx context.Context) {
	const shutdownDelay = 9 * time.Second // `docker stop` use 10s between SIGTERM and SIGKILL

	<-ctx.Done()
	time.Sleep(shutdownDelay)
	logger.With(zap.String(log.Version, ver)).Fatal("failed to graceful shutdown")
}
