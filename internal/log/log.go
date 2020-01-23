// Package log stores logged fields, and also provides helper methods for interaction with the logger.
package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Log field names.
const (
	Host       = "host"
	Port       = "port"
	Addr       = "addr"
	Remote     = "remote" // aligned IPv4:port "   192.168.0.42:1234 "
	Func       = "func"   // RPC method name, REST resource path
	HTTPMethod = "httpMethod"
	Error      = "error"
	HTTPStatus = "httpStatus"
	User       = "userID"
	API        = "api"
	GRPCCode   = "grpcCode"
	Version    = "version"
)

type loggerKey uint8

const logKey loggerKey = 1

// FromContext retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func FromContext(ctx context.Context) logrus.FieldLogger {
	val := ctx.Value(logKey)

	log, ok := val.(logrus.FieldLogger)
	if ok {
		return log
	}

	return logrus.New()
}

// SetContext puts a logger in context.
func SetContext(ctx context.Context, log logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, logKey, log)
}
