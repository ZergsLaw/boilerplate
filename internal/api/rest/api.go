// Package rest contains all methods and middleware for working web server.
package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sebest/xff"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/log"
	"go.uber.org/zap"
)

type (
	service struct {
		app app.App
	}

	config struct {
		host     string
		port     int
		basePath string
	}
	// Option for run server.
	Option func(*config)
)

// SetBasePath sets the base path to handlers.
// Default: /api/v1.
func SetBasePath(basePath string) Option {
	return func(c *config) {
		c.basePath = basePath
	}
}

// SetPort sets server port.
// Default: 8080.
func SetPort(port int) Option {
	return func(c *config) {
		c.port = port
	}
}

// SetHost sets server host.
// Default: localhost.
func SetHost(host string) Option {
	return func(c *config) {
		c.host = host
	}
}

func defaultConfig() *config {
	return &config{
		host:     "localhost",
		port:     8080,
		basePath: "",
	}
}

// New returns Swagger server configured to listen on the TCP network.
func New(application app.App, logger *zap.Logger, options ...Option) (*restapi.Server, error) {
	svc := &service{app: application}
	cfg := defaultConfig()

	for i := range options {
		options[i](cfg)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, fmt.Errorf("load embedded swagger spec: %w", err)
	}
	if cfg.basePath == "" {
		cfg.basePath = swaggerSpec.BasePath()
	}
	swaggerSpec.Spec().BasePath = cfg.basePath
	api := operations.NewServiceBoilerplateAPI(swaggerSpec)
	api.Logger = logger.Named("swagger").Sugar().Infof
	api.CookieKeyAuth = svc.cookieKeyAuth

	api.VerificationEmailHandler = operations.VerificationEmailHandlerFunc(svc.verificationEmail)
	api.VerificationUsernameHandler = operations.VerificationUsernameHandlerFunc(svc.verificationUsername)
	api.CreateUserHandler = operations.CreateUserHandlerFunc(svc.createUser)
	api.LoginHandler = operations.LoginHandlerFunc(svc.login)
	api.LogoutHandler = operations.LogoutHandlerFunc(svc.logout)
	api.GetUserHandler = operations.GetUserHandlerFunc(svc.getUser)
	api.DeleteUserHandler = operations.DeleteUserHandlerFunc(svc.deleteUser)
	api.UpdatePasswordHandler = operations.UpdatePasswordHandlerFunc(svc.updatePassword)
	api.UpdateUsernameHandler = operations.UpdateUsernameHandlerFunc(svc.updateUsername)
	api.UpdateEmailHandler = operations.UpdateEmailHandlerFunc(svc.updateEmail)
	api.GetUsersHandler = operations.GetUsersHandlerFunc(svc.getUsers)

	server := restapi.NewServer(api)
	server.Host = cfg.host
	server.Port = cfg.port

	// The middlewareFunc executes before anything.
	globalMiddlewares := func(handler http.Handler) http.Handler {
		xffmw, _ := xff.Default()
		createLog := createLogger(cfg.basePath, logger)
		accesslog := accessLog(cfg.basePath)
		redocOpts := middleware.RedocOpts{
			BasePath: cfg.basePath,
			SpecURL:  path.Join(cfg.basePath, "/swagger.json"),
		}

		return xffmw.Handler(createLog(recovery(accesslog(
			middleware.Spec(cfg.basePath, restapi.FlatSwaggerJSON,
				middleware.Redoc(redocOpts,
					handler))))))
	}

	server.SetHandler(globalMiddlewares(api.Serve(nil)))

	return server, nil
}

func fromRequest(r *http.Request, authUser *app.AuthUser) (context.Context, *zap.Logger, string) {
	ctx := r.Context()
	userID := app.UserID(0)
	if authUser != nil {
		userID = authUser.ID
	}

	logger := log.FromContext(ctx).With(zap.Int(log.User, int(userID)))
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ctx, logger, remoteIP
}
