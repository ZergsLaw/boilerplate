package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	zergrepo "github.com/ZergsLaw/zerg-repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"github.com/zergslaw/boilerplate/internal/api/rpc"
	"github.com/zergslaw/boilerplate/internal/api/web"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/auth"
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
	WebServerPort    = 8080
	GRPCServerPort   = 3000
	MetricServerPort = 9080
)

var (
	jwtKey = &cli.StringFlag{
		Name:     "jwt-key",
		Usage:    "jwt key for hashing auth",
		EnvVars:  []string{"JWT_KEY"},
		Required: true,
	}

	webHost = &cli.StringFlag{
		Name:    "web-host",
		Usage:   "web server host",
		EnvVars: []string{"SERVER_HOST"},
	}

	restPort = &cli.IntFlag{
		Name:     "web-port",
		Usage:    "web server port",
		EnvVars:  []string{"SERVER_PORT"},
		Required: true,
		Value:    WebServerPort,
	}

	metricHost = &cli.StringFlag{
		Name:    "metric-host",
		Usage:   "metric server host",
		EnvVars: []string{"METRIC_HOST"},
	}

	metricPort = &cli.IntFlag{
		Name:     "metric-port",
		Usage:    "metric server port",
		EnvVars:  []string{"METRIC_PORT"},
		Required: true,
		Value:    MetricServerPort,
	}

	gRPCHost = &cli.StringFlag{
		Name:    "gRPC-host",
		Usage:   "gRPC server host",
		EnvVars: []string{"GRPC_HOST"},
	}

	gRPCPort = &cli.IntFlag{
		Name:     "gRPC-port",
		Usage:    "gRPC server port",
		EnvVars:  []string{"GRPC_PORT"},
		Required: true,
		Value:    GRPCServerPort,
	}

	emailFrom = &cli.StringFlag{
		Name:     "email-from",
		Usage:    "email for notification",
		EnvVars:  []string{"EMAIL_FROM"},
		Required: true,
	}

	emailAPIKey = &cli.StringFlag{
		Name:     "email-api-key",
		Usage:    "set api key for send email",
		EnvVars:  []string{"EMAIL_API_KEY"},
		Required: true,
	}

	Serve = &cli.Command{
		Name:         "serve",
		Aliases:      []string{"s"},
		Usage:        "starts the service.",
		UsageText:    "Starts the service.",
		BashComplete: cli.DefaultAppComplete,
		Before:       beforeAction,
		Action:       serverAction,
		Flags: []cli.Flag{
			operation,
			dbName, dbUser, dbPass, dbHost, dbPort,
			jwtKey,
			webHost, restPort,
			metricHost, metricPort,
			gRPCHost, gRPCPort,
		},
	}
)

func beforeAction(c *cli.Context) error {
	if c.String(operation.Name) == "" {
		return nil
	}

	return migrateAction(c)
}

func serverAction(c *cli.Context) error {
	hostName, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("hostname: %w", err)
	}

	ctxConnect, cancelConnect := context.WithTimeout(c.Context, ConnectTimeout)
	defer cancelConnect()

	dbConn, err := zergrepo.ConnectByCfg(ctxConnect, "postgres", zergrepo.Config{
		Host:     c.String(dbHost.Name),
		Port:     c.Int(dbPort.Name),
		User:     c.String(dbUser.Name),
		Password: c.String(dbPass.Name),
		DBName:   c.String(dbName.Name),
		SSLMode:  zergrepo.DBSSLMode,
	})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	zp := repo.Connect(dbConn, log.FromContext(c.Context).Named("zergrepo").Sugar(), c.App.Name)
	r := repo.New(zp)

	emailClientConn, err := notification.Connect(c.String(emailAPIKey.Name))
	if err != nil {
		return fmt.Errorf("connect sendgrid: %w", err)
	}
	n := notification.New(emailClientConn, c.String(emailFrom.Name))

	pass := password.New()
	tokenizer := auth.New(c.String(jwtKey.Name))
	rc := recoverycode.New()
	application := app.New(app.Config{
		UserRepo: r, SessionRepo: r, CodeRepo: r, Wal: r,
		Password:     pass,
		Auth:         tokenizer,
		Notification: n,
		Code:         rc,
	})

	webAPIHost := host(c.String(webHost.Name), hostName)
	gRPCAPIHost := host(c.String(gRPCHost.Name), hostName)
	metricAPIHost := host(c.String(metricHost.Name), hostName)

	group, ctx := errgroup.WithContext(c.Context)
	services := []func() error{
		func() error { return webAPI(ctx, application, webAPIHost, c.Int(restPort.Name)) },
		func() error { return metricAPI(ctx, metricAPIHost, c.Int(metricPort.Name)) },
		func() error { return grpcAPI(ctx, application, gRPCAPIHost, c.Int(gRPCPort.Name)) },
		func() error { return startWAL(ctx, application) },
	}

	for _, service := range services {
		group.Go(service)
	}

	return group.Wait()
}

func host(host, defHost string) string {
	if host == "" {
		return defHost
	}

	return host
}

func webAPI(ctx context.Context, application app.App, host string, port int) error {
	logger := log.FromContext(ctx).Named("web")

	api, err := web.New(application,
		logger,
		web.SetHost(host),
		web.SetPort(port),
	)
	if err != nil {
		return fmt.Errorf("web new: %w", err)
	}

	errc := make(chan error, 1)
	go func() { errc <- api.Serve() }()
	logger.Info("server started", zap.String(log.Host, host), zap.Int(log.Port, port))

	select {
	case err = <-errc:
	case <-ctx.Done():
		err = api.Shutdown()
	}
	if err != nil {
		return fmt.Errorf("failed to serve Swagger protocol: %w", err)
	}

	logger.Info("shutdown server")
	return nil
}

func metricAPI(ctx context.Context, host string, port int) error {
	logger := log.FromContext(ctx).Named("prometheus")

	http.Handle("/metrics", promhttp.Handler())
	metricSrv := &http.Server{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
	}

	errc := make(chan error, 1)
	go func() { errc <- metricSrv.ListenAndServe() }()
	logger.Info("server started", zap.String(log.Host, host), zap.Int(log.Port, port))

	var err error
	select {
	case err = <-errc:
	case <-ctx.Done():
		err = metricSrv.Shutdown(context.Background())
	}
	if err != nil {
		return fmt.Errorf("failed to serve prometheus metrics: %w", err)
	}

	logger.Info("shutdown server")
	return nil
}

func grpcAPI(ctx context.Context, application app.UserApp, host string, port int) error {
	logger := log.FromContext(ctx).Named("gRPC")

	api := rpc.New(application, logger)
	ln, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	errc := make(chan error, 1)
	go func() { errc <- api.Serve(ln) }()

	logger.Info("server started", zap.String(log.Host, host), zap.Int(log.Port, port))

	select {
	case err = <-errc:
	case <-ctx.Done():
		api.GracefulStop()
	}
	if err != nil {
		return fmt.Errorf("failed to serve prometheus metrics: %w", err)
	}

	logger.Info("shutdown server")
	return nil
}

func startWAL(ctx context.Context, application app.WALApplication) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return application.StartWALNotification(ctx)
	})

	return group.Wait()
}
