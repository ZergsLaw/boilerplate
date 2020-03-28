package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/api/rest"
	"github.com/zergslaw/boilerplate/internal/api/rpc"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/auth"
	"github.com/zergslaw/boilerplate/internal/flag"
	"github.com/zergslaw/boilerplate/internal/log"
	"github.com/zergslaw/boilerplate/internal/notification"
	"github.com/zergslaw/boilerplate/internal/password"
	"github.com/zergslaw/boilerplate/internal/recoverycode"
	"github.com/zergslaw/boilerplate/internal/repo"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// Default values.
const (
	RestServerPort   = 8080
	GRPCServerPort   = 3000
	MetricServerPort = 9080

	RabbitHost  = "localhost"
	RabbitUser  = "rabbit"
	RabbitPass  = "rabbit"
	RabbitPort  = 5672
	RabbitQueue = "notification"
)

// nolint:gochecknoglobals,gocritic
var (
	migrateFlag = flag.NewStrFlag("migrate", "goose migrate when you start the service",
		flag.StrEnv("MIGRATE"))

	jwtKey = flag.NewStrFlag("jwt-key", "jwt key for hashing auth",
		flag.StrRequired(), flag.StrAliases("JWT_KEY"))

	restHost = flag.NewStrFlag("rest-host", "rest server host",
		flag.StrRequired(), flag.StrAliases("SERVER_HOST"), flag.StrDefault(host))
	restPort = flag.NewIntFlag("rest-port", "rest server port",
		flag.IntRequired(), flag.IntAliases("SERVER_PORT"), flag.IntDefault(RestServerPort))

	metricHost = flag.NewStrFlag("metric-host", "serve prometheus metrics on host",
		flag.StrRequired(), flag.StrAliases("METRIC_HOST"), flag.StrDefault(host))
	metricPort = flag.NewIntFlag("metric-port", "serve prometheus metrics on port",
		flag.IntRequired(), flag.IntAliases("METRIC_PORT"), flag.IntDefault(MetricServerPort))

	grpcHost = flag.NewStrFlag("gRPC-host", "serve internal gRPC API on host",
		flag.StrRequired(), flag.StrAliases("GRPC_HOST"), flag.StrDefault(host))
	grpcPort = flag.NewIntFlag("gRPC-port", "serve internal gRPC API on port",
		flag.IntRequired(), flag.IntAliases("GRPC_PORT"), flag.IntDefault(GRPCServerPort))

	rabbitUser = flag.NewStrFlag("rabbit-user", "rabbit user",
		flag.StrRequired(), flag.StrEnv("RABBIT_USER"), flag.StrDefault(RabbitUser))
	rabbitPass = flag.NewStrFlag("rabbit-pass", "rabbit pass",
		flag.StrRequired(), flag.StrEnv("RABBIT_PASS"), flag.StrDefault(RabbitPass))
	rabbitHost = flag.NewStrFlag("rabbit-host", "rabbit host",
		flag.StrRequired(), flag.StrEnv("RABBIT_HOST"), flag.StrDefault(RabbitHost))
	rabbitPort = flag.NewIntFlag("rabbit-port", "rabbit port",
		flag.IntRequired(), flag.IntEnv("RABBIT_PORT"), flag.IntDefault(RabbitPort))
	queueName = flag.NewStrFlag("queue", "queue name for sending notification",
		flag.StrRequired(), flag.StrEnv("RABBIT_PORT"), flag.StrDefault(RabbitQueue))

	serve = &cli.Command{
		Name:         "serve",
		Aliases:      []string{"s"},
		Usage:        "starts the service.",
		UsageText:    "Starts the service.",
		BashComplete: cli.DefaultAppComplete,
		Before:       beforeAction,
		Action:       serverAction,
		Flags: []cli.Flag{
			migrateFlag,
			dbName, dbUser, dbPass, dbHost, dbPort, migrateDir,
			jwtKey,
			restHost, restPort,
			metricHost, metricPort,
			grpcHost, grpcPort,
			rabbitUser, rabbitHost, rabbitPass, rabbitPort,
		},
	}
)

func beforeAction(c *cli.Context) error {
	if c.String(migrateFlag.Name) == "" {
		return nil
	}

	return goose(c.Context, c.String(migrateDir.Name), c.String(migrateFlag.Name),
		repo.Name(c.String(dbName.Name)),
		repo.User(c.String(dbUser.Name)),
		repo.Pass(c.String(dbPass.Name)),
		repo.Host(c.String(dbHost.Name)),
		repo.Port(c.Int(dbPort.Name)),
	)
}

func serverAction(c *cli.Context) error {
	ctxConnect, cancelConnect := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancelConnect()

	dbConn, err := repo.Connect(ctxConnect,
		repo.Name(c.String(dbName.Name)),
		repo.User(c.String(dbUser.Name)),
		repo.Pass(c.String(dbPass.Name)),
		repo.Host(c.String(dbHost.Name)),
		repo.Port(c.Int(dbPort.Name)),
	)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	rabbitConn, err := notification.Connect(notification.Config{
		User: c.String(rabbitUser.Name),
		Pass: c.String(rabbitPass.Name),
		Host: c.String(rabbitHost.Name),
		Port: c.Int(rabbitPort.Name),
	})
	if err != nil {
		return fmt.Errorf("connect rabbit mq: %w", err)
	}

	ch, err := rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("get rabbit channel: %w", err)
	}

	_, err = ch.QueueDeclare(c.String(queueName.Name), false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	r := repo.New(dbConn)
	n := notification.New(ch)
	pass := password.New()
	tokenizer := auth.New(c.String(jwtKey.Name))
	rc := recoverycode.New()
	application := app.New(r, r, r, pass, tokenizer, r, n, rc)

	group, ctx := errgroup.WithContext(c.Context)
	services := []func() error{
		func() error { return swaggerAPI(ctx, application, c.String(restHost.Name), c.Int(restPort.Name)) },
		func() error { return metricAPI(ctx, c.String(metricHost.Name), c.Int(metricPort.Name)) },
		func() error { return grpcAPI(ctx, application, c.String(grpcHost.Name), c.Int(grpcPort.Name)) },
		func() error { return startWAL(ctx, application) },
	}

	for _, service := range services {
		group.Go(service)
	}

	return group.Wait()
}

func swaggerAPI(ctx context.Context, application app.App, host string, port int) error {
	restLogger := logger.Named("rest")

	api, err := rest.New(application,
		restLogger,
		rest.SetHost(host),
		rest.SetPort(port),
	)
	if err != nil {
		return fmt.Errorf("rest new: %w", err)
	}

	errc := make(chan error, 1)
	go func() { errc <- api.Serve() }()

	restLogger.With(
		zap.String(log.Host, host),
		zap.Int(log.Port, port),
	).Info("serve server")

	select {
	case err = <-errc:
	case <-ctx.Done():
		err = api.Shutdown()
	}
	if err != nil {
		return fmt.Errorf("failed to serve Swagger protocol: %w", err)
	}

	restLogger.Info("shutdown server")
	return nil
}

func metricAPI(ctx context.Context, host string, port int) error {
	http.Handle("/metrics", promhttp.Handler())
	metricSrv := &http.Server{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
	}

	logger.Named("prometheus").With(
		zap.String(log.Host, host),
		zap.Int(log.Port, port),
	).Info("serve metrics")

	errc := make(chan error, 1)
	go func() { errc <- metricSrv.ListenAndServe() }()

	var err error
	select {
	case err = <-errc:
	case <-ctx.Done():
		err = metricSrv.Shutdown(context.Background())
	}
	if err != nil {
		return fmt.Errorf("failed to serve Prometheus metrics: %w", err)
	}

	logger.Info("shutdown Prometheus metrics")
	return nil
}

func grpcAPI(ctx context.Context, application app.UserApp, host string, port int) error {
	gRPCLogger := logger.Named("gRPC")

	api := rpc.New(application, gRPCLogger)
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	gRPCLogger.With(
		zap.String(log.Host, host),
		zap.Int(log.Port, port),
	).Info("serve server")

	go func() { <-ctx.Done(); api.GracefulStop() }()

	err = api.Serve(ln)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	gRPCLogger.Info("shutdown server")
	return nil
}

func startWAL(ctx context.Context, application app.WALApplication) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return application.StartWALNotification(ctx)
	})

	return group.Wait()
}
