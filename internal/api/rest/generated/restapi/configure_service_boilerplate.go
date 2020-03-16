// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi/operations"
	"github.com/zergslaw/boilerplate/internal/app"
)

//go:generate swagger generate server --target ../../generated --name ServiceBoilerplate --spec ../../swagger.yml --principal app.AuthUser --exclude-main

func configureFlags(api *operations.ServiceBoilerplateAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ServiceBoilerplateAPI) http.Handler {
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
	if api.CookieKeyAuth == nil {
		api.CookieKeyAuth = func(token string) (*app.AuthUser, error) {
			return nil, errors.NotImplemented("api key auth (cookieKey) Cookie from header param [Cookie] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.CreateRecoveryCodeHandler == nil {
		api.CreateRecoveryCodeHandler = operations.CreateRecoveryCodeHandlerFunc(func(params operations.CreateRecoveryCodeParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CreateRecoveryCode has not yet been implemented")
		})
	}
	if api.CreateUserHandler == nil {
		api.CreateUserHandler = operations.CreateUserHandlerFunc(func(params operations.CreateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CreateUser has not yet been implemented")
		})
	}
	if api.DeleteUserHandler == nil {
		api.DeleteUserHandler = operations.DeleteUserHandlerFunc(func(params operations.DeleteUserParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteUser has not yet been implemented")
		})
	}
	if api.GetUserHandler == nil {
		api.GetUserHandler = operations.GetUserHandlerFunc(func(params operations.GetUserParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUser has not yet been implemented")
		})
	}
	if api.GetUsersHandler == nil {
		api.GetUsersHandler = operations.GetUsersHandlerFunc(func(params operations.GetUsersParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUsers has not yet been implemented")
		})
	}
	if api.LoginHandler == nil {
		api.LoginHandler = operations.LoginHandlerFunc(func(params operations.LoginParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.Login has not yet been implemented")
		})
	}
	if api.LogoutHandler == nil {
		api.LogoutHandler = operations.LogoutHandlerFunc(func(params operations.LogoutParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.Logout has not yet been implemented")
		})
	}
	if api.RecoveryPasswordHandler == nil {
		api.RecoveryPasswordHandler = operations.RecoveryPasswordHandlerFunc(func(params operations.RecoveryPasswordParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.RecoveryPassword has not yet been implemented")
		})
	}
	if api.UpdateEmailHandler == nil {
		api.UpdateEmailHandler = operations.UpdateEmailHandlerFunc(func(params operations.UpdateEmailParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.UpdateEmail has not yet been implemented")
		})
	}
	if api.UpdatePasswordHandler == nil {
		api.UpdatePasswordHandler = operations.UpdatePasswordHandlerFunc(func(params operations.UpdatePasswordParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.UpdatePassword has not yet been implemented")
		})
	}
	if api.UpdateUsernameHandler == nil {
		api.UpdateUsernameHandler = operations.UpdateUsernameHandlerFunc(func(params operations.UpdateUsernameParams, principal *app.AuthUser) middleware.Responder {
			return middleware.NotImplemented("operation operations.UpdateUsername has not yet been implemented")
		})
	}
	if api.VerificationEmailHandler == nil {
		api.VerificationEmailHandler = operations.VerificationEmailHandlerFunc(func(params operations.VerificationEmailParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.VerificationEmail has not yet been implemented")
		})
	}
	if api.VerificationUsernameHandler == nil {
		api.VerificationUsernameHandler = operations.VerificationUsernameHandlerFunc(func(params operations.VerificationUsernameParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.VerificationUsername has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

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
