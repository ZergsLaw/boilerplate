package rpc_test

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/zergslaw/boilerplate/internal/api/rpc/pb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/boilerplate/internal/api/rpc"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/metrics"
	"github.com/zergslaw/boilerplate/internal/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	metrics.InitMetrics()

	os.Exit(m.Run())
}

var (
	errAny  = errors.New("any err")
	ctx     = context.Background()
	token   = "token"
	rpcUser = pb.User{
		Id:       1,
		Username: "username",
		Email:    "email@email.com",
	}
	appUser = app.AuthUser{
		User: app.User{
			ID:    app.UserID(rpcUser.Id),
			Email: rpcUser.Email,
			Name:  rpcUser.Username,
		},
	}
)

func testNew(t *testing.T) (pb.UsersClient, *mock.MockApp, func()) {
	t.Helper()

	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.FatalLevel))
	assert.Nil(t, err)

	ctrl := gomock.NewController(t)
	mockApp := mock.NewMockApp(ctrl)
	server := rpc.New(mockApp, logger)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)

	go func() {
		err := server.Serve(ln)
		assert.Nil(t, err)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.DialContext(ctx, ln.Addr().String(),
		grpc.WithInsecure(), // TODO Add TLS and remove this.
		grpc.WithBlock(),
	)
	assert.Nil(t, err)

	shutdown := func() {
		t.Helper()
		assert.Nil(t, conn.Close())
		ctrl.Finish()
		server.GracefulStop()
		cancel()
	}

	return pb.NewUsersClient(conn), mockApp, shutdown
}
