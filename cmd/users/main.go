package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-openapi/loads"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zergslaw/users/internal/api/rest"
	"github.com/zergslaw/users/internal/api/rest/generated/restapi"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/auth"
	"github.com/zergslaw/users/internal/config"
	"github.com/zergslaw/users/internal/db"
	"github.com/zergslaw/users/internal/password"
	"github.com/zergslaw/users/migration"
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

	ServerHost = "localhost"
	ServerPort = 8080
)

// nolint:gochecknoglobals
var (
	log = logrus.New()
	exe = strings.TrimSuffix(path.Base(os.Args[0]), ".test")
	ver string

	cfg struct {
		rest struct {
			jwtKey   string
			host     string
			port     int
			basePath string
			migrate  string
		}

		db struct {
			name     string
			user     string
			pass     string
			host     string
			port     int
			gooseDir string
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
				log.Fatal(err)
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
					log.Fatal(err)
				}
			}

			err := serve()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func goose(command string) error {
	return migration.Run(context.Background(), cfg.db.gooseDir, command,
		db.Name(cfg.db.name),
		db.User(cfg.db.user),
		db.Pass(cfg.db.pass),
		db.Host(cfg.db.host),
		db.Port(cfg.db.port),
	)
}

func serve() error {
	log.Info("started ", "version ", ver)
	defer log.Info("finished ", "version ", ver)

	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer cancel()

	dbConn, err := db.Connect(ctx,
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
	server, err := rest.New(application,
		rest.SetHost(cfg.rest.host),
		rest.SetPort(cfg.rest.port),
	)
	if err != nil {
		return fmt.Errorf("rest build: %w", err)
	}

	return server.Serve()
}

func init() { // nolint:gochecknoinits
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(fmt.Errorf("load embedded swagger spec: %w", err))
	}

	ver = fmt.Sprintf("%s  %s", swaggerSpec.Spec().Info.Version, runtime.Version())

	setFlagsDB(gooseCmd)
	setFlagsDB(serveCmd)
	setFlagsRestServer(serveCmd)
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

func setFlagsRestServer(cmd *cobra.Command) {
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(fmt.Errorf("get hostname: %w", err))
	}

	cmd.Flags().StringVar(&cfg.rest.jwtKey, "jwt-key", config.Env("JWT_KEY"), "jwt key for hashing auth")
	cmd.Flags().StringVar(&cfg.rest.host, "rest-host", config.EnvOrDef("SERVER_HOST", host), "rest host")
	cmd.Flags().IntVar(&cfg.rest.port, "rest-port", config.IntEnvOrDef("SERVER_PORT", ServerPort), "rest port")
	cmd.Flags().StringVar(&cfg.rest.migrate, "migrate", config.Env("MIGRATE"), "goose migrate when you start the rest")

	err = cmd.MarkFlagRequired("jwt-key")
	if err != nil {
		log.Fatal(fmt.Errorf("mark required flag: %w", err))
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
