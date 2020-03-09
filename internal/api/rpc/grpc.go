// Package rpc contains all methods and middleware for working gRPC server.
package rpc

import (
	"context"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/zergslaw/boilerplate/internal/app"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type users interface {
	// UserByAuthToken is documented in app.App interface.
	UserByAuthToken(ctx context.Context, token app.AuthToken) (*app.AuthUser, error)
}

type service struct {
	app users
}

// New returns gRPC server configured to listen on the TCP network.
func New(application users, logger *zap.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    50 * time.Second,
			Timeout: 10 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			prometheus.UnaryServerInterceptor,
			MakeUnaryServerLogger(logger),
			UnaryServerRecover,
			UnaryServerAccessLog,
		)),
	)

	RegisterUsersServer(server, &service{app: application})

	prometheus.Register(server)

	return server
}
