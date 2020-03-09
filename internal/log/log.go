// Package log stores logged fields, and also provides helper methods for interaction with the logger.
package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
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
func FromContext(ctx context.Context) *zap.Logger {
	val := ctx.Value(logKey)

	log, ok := val.(*zap.Logger)
	if ok {
		return log
	}

	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("init new logger: %w", err))
	}

	return log
}

// SetContext puts a logger in context.
func SetContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}
