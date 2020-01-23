package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/zergslaw/users/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type service struct {
	app app.Users
}

func NewServer(application app.Users) *grpc.Server {
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    50 * time.Second,
			Timeout: 10 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			MakeUnaryServerLogger,
			UnaryServerRecover,
			UnaryServerAccessLog,
		)),
	)

	RegisterUsersServer(server, &service{app: application})

	grpc_prometheus.Register(server)

	return server
}
