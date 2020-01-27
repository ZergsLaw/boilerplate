package rpc_test

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/users/internal/api/rpc"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/metrics"
	"github.com/zergslaw/users/internal/mock"
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
	rpcUser = rpc.User{
		Id:       1,
		Username: "username",
		Email:    "email@email.com",
	}
	appUser = app.AuthUser{
		User: app.User{
			ID:       app.UserID(rpcUser.Id),
			Email:    rpcUser.Email,
			Username: rpcUser.Username,
		},
	}
)

func testNew(t *testing.T) (rpc.UsersClient, *mock.App, func()) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockApp := mock.NewApp(ctrl)
	server := rpc.New(mockApp)

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

	return rpc.NewUsersClient(conn), mockApp, shutdown
}
