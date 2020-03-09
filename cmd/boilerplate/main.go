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

	"github.com/go-openapi/loads"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/api/rest"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi"
	"github.com/zergslaw/boilerplate/internal/log"
	"github.com/zergslaw/boilerplate/internal/repo"
	"go.uber.org/zap"
)

// nolint:gochecknoglobals,gocritic
var (
	logger *zap.Logger
	ver    string
	host   string
	appl   = &cli.App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "Boilerplate application.",
		BashComplete: cli.DefaultAppComplete,
		Writer:       os.Stdout,
		Commands:     []*cli.Command{version, migrate, serve},
	}
)

func initDefaultData() error {
	var err error
	host, err = os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname: %w", err)
	}

	logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return fmt.Errorf("load embedded swagger spec: %w", err)
	}
	ver = swaggerSpec.Spec().Info.Version

	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(appl.Name, "_")
	repo.InitMetrics(namespace)
	rest.InitMetrics(namespace, restapi.FlatSwaggerJSON)

	return nil
}

func main() {
	err := initDefaultData()
	if err != nil {
		logger.Fatal("init service", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() { <-signals; cancel() }()
	go forceShutdown(ctx)

	err = appl.RunContext(ctx, os.Args)
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
