// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/zergslaw/users/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/users/internal/app"
)

//go:generate swagger generate server --target ../../generated --name ServiceUser --spec ../../swagger.yml --principal app.AuthUser --exclude-main --strict

func configureFlags(api *operations.ServiceUserAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ServiceUserAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Cookie" header is set
	api.CookieKeyAuth = func(token string) (*app.AuthUser, error) {
		return nil, errors.NotImplemented("api key auth (cookieKey) Cookie from header param [Cookie] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.CreateUserHandler = operations.CreateUserHandlerFunc(func(params operations.CreateUserParams) operations.CreateUserResponder {
		return operations.CreateUserNotImplemented()
	})
	api.DeleteUserHandler = operations.DeleteUserHandlerFunc(func(params operations.DeleteUserParams, principal *app.AuthUser) operations.DeleteUserResponder {
		return operations.DeleteUserNotImplemented()
	})
	api.GetUserHandler = operations.GetUserHandlerFunc(func(params operations.GetUserParams) operations.GetUserResponder {
		return operations.GetUserNotImplemented()
	})
	api.GetUsersHandler = operations.GetUsersHandlerFunc(func(params operations.GetUsersParams, principal *app.AuthUser) operations.GetUsersResponder {
		return operations.GetUsersNotImplemented()
	})
	api.LoginHandler = operations.LoginHandlerFunc(func(params operations.LoginParams) operations.LoginResponder {
		return operations.LoginNotImplemented()
	})
	api.LogoutHandler = operations.LogoutHandlerFunc(func(params operations.LogoutParams, principal *app.AuthUser) operations.LogoutResponder {
		return operations.LogoutNotImplemented()
	})
	api.UpdateEmailHandler = operations.UpdateEmailHandlerFunc(func(params operations.UpdateEmailParams, principal *app.AuthUser) operations.UpdateEmailResponder {
		return operations.UpdateEmailNotImplemented()
	})
	api.UpdatePasswordHandler = operations.UpdatePasswordHandlerFunc(func(params operations.UpdatePasswordParams, principal *app.AuthUser) operations.UpdatePasswordResponder {
		return operations.UpdatePasswordNotImplemented()
	})
	api.UpdateUsernameHandler = operations.UpdateUsernameHandlerFunc(func(params operations.UpdateUsernameParams, principal *app.AuthUser) operations.UpdateUsernameResponder {
		return operations.UpdateUsernameNotImplemented()
	})
	api.VerificationEmailHandler = operations.VerificationEmailHandlerFunc(func(params operations.VerificationEmailParams) operations.VerificationEmailResponder {
		return operations.VerificationEmailNotImplemented()
	})
	api.VerificationUsernameHandler = operations.VerificationUsernameHandlerFunc(func(params operations.VerificationUsernameParams) operations.VerificationUsernameResponder {
		return operations.VerificationUsernameNotImplemented()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
