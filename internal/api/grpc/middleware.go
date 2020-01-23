package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zergslaw/users/internal/log"
	"github.com/zergslaw/users/internal/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"path"
)

// MakeUnaryServerLogger returns a new unary server interceptor that contains request logger.
func MakeUnaryServerLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	logger := newLogger(ctx, info.FullMethod)

	ctx = log.SetContext(ctx, logger)
	return handler(ctx, req)
}

// UnaryServerRecover returns a new unary server interceptor that recover and logs panic.
func UnaryServerRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			metrics.PanicsTotal.Inc()
			logger := log.FromContext(ctx)
			logger.WithFields(logrus.Fields{
				log.GRPCCode: codes.Internal,
				log.Error:    "panic",
			}).Error(p)

			err = status.Errorf(codes.Internal, "%v", p)
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

func newLogger(ctx context.Context, fullMethod string) logrus.FieldLogger {
	remoteAddr := ""
	if p, ok := peer.FromContext(ctx); ok {
		remoteAddr = p.Addr.String()
	}

	logger := logrus.New().WithFields(logrus.Fields{
		log.Remote:   remoteAddr,
		log.Func:     path.Base(fullMethod),
		log.API:      "grpc",
		log.GRPCCode: "",
	})

	return logger
}

func logHandler(logger logrus.FieldLogger, err error) {
	s := status.Convert(err)
	code, msg := s.Code(), s.Message()
	switch code {
	case codes.OK, codes.Canceled, codes.NotFound:
		logger.Info("handled", log.GRPCCode, code)
	default:
		logger.Error("failed to handle", log.GRPCCode, code, "err", msg)
	}
}
