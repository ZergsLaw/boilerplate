package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-openapi/loads"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zergslaw/users/internal/api/rest"
	"github.com/zergslaw/users/internal/api/rest/generated/restapi"
	"github.com/zergslaw/users/internal/api/rpc"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/auth"
	"github.com/zergslaw/users/internal/config"
	"github.com/zergslaw/users/internal/db"
	"github.com/zergslaw/users/internal/log"
	"github.com/zergslaw/users/internal/password"
	"github.com/zergslaw/users/migration"
	"golang.org/x/sync/errgroup"
)

// Default value.
const (
	DBName = "postgres"
	DBUser = "postgres"
	DBPass = "postgres"
	DBHost = "localhost"
	DBPort = 5432

	DirMigrate = "migration"

	ConnectTimeout = time.Second * 5

	RestServerPort   = 8080
	GRPCServerPort   = 3000
	MetricServerPort = 9080
)

// nolint:gochecknoglobals
var (
	logger = logrus.New()
	exe    = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	host   string
	ver    string

	cfg struct {
		rest struct {
			jwtKey   string
			host     string
			port     int
			basePath string
			migrate  string
		}

		grpc struct {
			host string
			port int
		}

		db struct {
			name     string
			user     string
			pass     string
			host     string
			port     int
			gooseDir string
		}

		metric struct {
			host string
			port int
		}
	}

	rootCmd = &cobra.Command{
		Use:     exe,
		Short:   "Microservice user.",
		Version: ver,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the service version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(exe, "version", ver)
		},
	}

	gooseCmd = &cobra.Command{
		Use:   "goose",
		Short: "Migrate database schema",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command := strings.Join(args, " ")

			err := goose(command)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the service",
		Run: func(cmd *cobra.Command, args []string) {

			if cfg.rest.migrate != "" {
				err := goose(cfg.rest.migrate)
				if err != nil {
					logger.Fatal(err)
				}
			}

			err := serve()
			if err != nil {
				logger.Fatal(err)
			}
		},
	}
)

func init() { // nolint:gochecknoinits
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logger.Fatal(fmt.Errorf("load embedded swagger spec: %w", err))
	}

	host, err = os.Hostname()
	if err != nil {
		logger.Fatal(fmt.Errorf("get hostname: %w", err))
	}

	ver = fmt.Sprintf("%s  %s", swaggerSpec.Spec().Info.Version, runtime.Version())

	setFlagsDB(gooseCmd)
	setFlagsDB(serveCmd)
	setFlagsServe(serveCmd)
	rootCmd.AddCommand(versionCmd, gooseCmd, serveCmd)

	namespace := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(exe, "_")
	db.InitMetrics(namespace)
	rest.InitMetrics(namespace, restapi.FlatSwaggerJSON)
}

func setFlagsDB(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.db.gooseDir, "dir-migrate", config.EnvOrDef("DB_MIGRATE", DirMigrate), "goose migrations dir")
	cmd.Flags().StringVar(&cfg.db.name, "db-name", config.EnvOrDef("DB_NAME", DBName), "database name")
	cmd.Flags().StringVar(&cfg.db.user, "db-user", config.EnvOrDef("DB_USER", DBUser), "database user")
	cmd.Flags().StringVar(&cfg.db.pass, "db-pass", config.EnvOrDef("DB_PASS", DBPass), "database pass")
	cmd.Flags().StringVar(&cfg.db.host, "db-host", config.EnvOrDef("DB_HOST", DBHost), "database host")
	cmd.Flags().IntVar(&cfg.db.port, "db-port", config.IntEnvOrDef("DB_PORT", DBPort), "database port")
}

func setFlagsServe(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.rest.jwtKey, "jwt-key", config.Env("JWT_KEY"), "jwt key for hashing auth")
	cmd.Flags().StringVar(&cfg.rest.host, "rest-host", config.EnvOrDef("SERVER_HOST", host), "rest host")
	cmd.Flags().IntVar(&cfg.rest.port, "rest-port", config.IntEnvOrDef("SERVER_PORT", RestServerPort), "rest port")
	cmd.Flags().StringVar(&cfg.rest.migrate, "migrate", config.Env("MIGRATE"), "goose migrate when you start the rest")

	cmd.Flags().StringVar(&cfg.metric.host, "metric-host", config.EnvOrDef("METRIC_HOST", host), "serve prometheus metrics on host")
	cmd.Flags().IntVar(&cfg.metric.port, "metric-port", config.IntEnvOrDef("METRIC_PORT", MetricServerPort), "serve prometheus metrics on port")

	cmd.Flags().StringVar(&cfg.grpc.host, "gRPC-host", config.EnvOrDef("GRPC_HOST", host), "serve internal gRPC API on host")
	cmd.Flags().IntVar(&cfg.grpc.port, "gRPC-port", config.IntEnvOrDef("GRPC_PORT", GRPCServerPort), "serve internal gRPC API on port")

	err := cmd.MarkFlagRequired("jwt-key")
	if err != nil {
		logger.Fatal(fmt.Errorf("mark required flag: %w", err))
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Fatal(err)
	}
}

func goose(command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	return migration.Run(ctx, cfg.db.gooseDir, command,
		db.Name(cfg.db.name),
		db.User(cfg.db.user),
		db.Pass(cfg.db.pass),
		db.Host(cfg.db.host),
		db.Port(cfg.db.port),
	)
}

func forceShutdown(ctx context.Context) {
	const shutdownDelay = 9 * time.Second // `docker stop` use 10s between SIGTERM and SIGKILL

	<-ctx.Done()
	time.Sleep(shutdownDelay)
	logger.WithField(log.Version, ver).Fatalln("failed to graceful shutdown")
}

func serve() error {
	logger.Info("started ", "version ", ver)
	defer logger.Info("finished ", "version ", ver)

	ctxConnect, cancelConnect := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancelConnect()

	dbConn, err := db.Connect(ctxConnect,
		db.Name(cfg.db.name),
		db.User(cfg.db.user),
		db.Pass(cfg.db.pass),
		db.Host(cfg.db.host),
		db.Port(cfg.db.port),
	)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}

	repo := db.New(dbConn)
	pass := password.New()
	tokenizer := auth.New(cfg.rest.jwtKey)
	application := app.New(repo, pass, tokenizer)

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() { <-signals; cancel() }()
	go forceShutdown(ctx)

	group, ctx := errgroup.WithContext(ctx)
	services := []func() error{
		func() error { return swaggerAPI(ctx, application) },
		func() error { return metricAPI(ctx) },
		func() error { return grpcAPI(ctx, application) },
	}

	for _, service := range services {
		group.Go(service)
	}

	return group.Wait()
}

func metricAPI(ctx context.Context) error {
	http.Handle("/metrics", promhttp.Handler())
	metricSrv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.metric.host, cfg.metric.port),
	}
	logger.WithFields(logrus.Fields{
		log.Host: cfg.metric.host,
		log.Port: cfg.metric.port,
	}).Info("serve Prometheus metrics")

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

func swaggerAPI(ctx context.Context, application app.App) error {
	api, err := rest.New(application,
		rest.SetHost(cfg.rest.host),
		rest.SetPort(cfg.rest.port),
	)
	if err != nil {
		return fmt.Errorf("rest new: %w", err)
	}

	logger.WithFields(logrus.Fields{
		log.Host: cfg.rest.host,
		log.Port: cfg.rest.port,
	}).Info("serve Swagger protocol")

	errc := make(chan error, 1)
	go func() { errc <- api.Serve() }()

	select {
	case err = <-errc:
	case <-ctx.Done():
		err = api.Shutdown()
	}
	if err != nil {
		return fmt.Errorf("failed to serve Swagger protocol: %w", err)
	}

	logger.Info("shutdown Swagger protocol")
	return nil
}

func grpcAPI(ctx context.Context, application app.App) error {
	api := rpc.New(application)
	ln, err := net.Listen("tcp", net.JoinHostPort(cfg.grpc.host, strconv.Itoa(cfg.grpc.port)))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	logger.WithFields(logrus.Fields{
		log.Host: cfg.grpc.host,
		log.Port: cfg.grpc.port,
	}).Info("serve gRPC protocol")

	go func() { <-ctx.Done(); api.GracefulStop() }()

	err = api.Serve(ln)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC protoco: %w", err)
	}

	logger.Info("shutdown gRPC protocol")
	return nil
}
