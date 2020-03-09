package rpc

import (
	"context"
	"path"

	"github.com/zergslaw/boilerplate/internal/log"
	"github.com/zergslaw/boilerplate/internal/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// MakeUnaryServerLogger returns a new unary server interceptor that contains request logger.
func MakeUnaryServerLogger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		remoteAddr := ""
		if p, ok := peer.FromContext(ctx); ok {
			remoteAddr = p.Addr.String()
		}

		logger.With(
			zap.String(log.Remote, remoteAddr),
			zap.String(log.Func, path.Base(path.Base(info.FullMethod))),
			zap.String(log.GRPCCode, ""),
		)

		ctx = log.SetContext(ctx, logger)
		return handler(ctx, req)
	}
}

// UnaryServerRecover returns a new unary server interceptor that recover and logs panic.
func UnaryServerRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			metrics.PanicsTotal.Inc()
			logger := log.FromContext(ctx)

			err = status.Errorf(codes.Internal, "%v", p)
			logger.Error("gRPC server panic",
				zap.String(log.GRPCCode, codes.Internal.String()),
				zap.Error(err),
			)
		}
	}()

	return handler(ctx, req)
}

// UnaryServerAccessLog returns a new unary server interceptor that logs request status.
func UnaryServerAccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	resp, err := handler(ctx, req)
	logger := log.FromContext(ctx)
	logHandler(logger, err)
	return resp, err
}

func logHandler(logger *zap.Logger, err error) {
	s := status.Convert(err)
	code, msg := s.Code(), s.Message()
	switch code {
	case codes.OK, codes.Canceled, codes.NotFound:
		logger.Info("handled", zap.String(log.GRPCCode, code.String()))
	default:
		logger.Error("failed to handle", zap.String(log.GRPCCode, code.String()), zap.String("msg", msg))
	}
}
