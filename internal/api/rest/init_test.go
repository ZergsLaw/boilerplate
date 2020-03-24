package rest_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/boilerplate/internal/api/rest"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/client"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/client/operations"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/models"
	"github.com/zergslaw/boilerplate/internal/api/rest/generated/restapi"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/metrics"
	"github.com/zergslaw/boilerplate/internal/mock"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	metrics.InitMetrics()
	rest.InitMetrics("test", restapi.FlatSwaggerJSON)

	os.Exit(m.Run())
}

var (
	errAny = errors.New("any error")

	notExistEmail    = "notExist@email.com"
	email            = "exist@email.com"
	notExistUsername = "notExistUsername"
	username         = "username"
	password         = "password"

	authToken app.AuthToken = "token"
	user                    = app.User{
		ID:        1,
		Email:     email,
		Username:  username,
		PassHash:  []byte(password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	origin = app.Origin{
		IP:        net.ParseIP("127.0.0.1"),
		UserAgent: "Go-http-client/1.1",
	}
	session = app.Session{
		Origin:  origin,
		ID:      1,
		TokenID: "tokenID",
	}

	authUser = app.AuthUser{
		User:    user,
		Session: session,
	}

	sessUser   = "sessUser"
	apiKeyAuth = httptransport.APIKeyAuth("Cookie", "header", "__Secure-authKey="+sessUser)
	restUser   = rest.User(&user)
)

func testNewServer(t *testing.T) (string, func(), *mock.App, *client.ServiceBoilerplate) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockApp := mock.NewApp(ctrl)
	mockApp.EXPECT().UserByAuthToken(gomock.Any(), app.AuthToken(sessUser)).
		Return(&authUser, nil).AnyTimes()

	log, err := zap.NewDevelopment(zap.AddStacktrace(zap.FatalLevel))
	assert.NoError(t, err)

	randomPort := rest.SetPort(0)
	server, err := rest.New(mockApp, log, randomPort)
	assert.NoError(t, err, "NewServer")
	assert.NoError(t, server.Listen(), "server.Listen")

	errc := make(chan error, 1)
	go func() { errc <- server.Serve() }()

	shutdown := func() {
		t.Helper()
		assert.Nil(t, server.Shutdown(), "server.Shutdown")
		assert.Nil(t, <-errc, "server.Serve")
		ctrl.Finish()
	}

	url := fmt.Sprintf("%s:%d", client.DefaultHost, server.Port)

	transport := httptransport.New(url, client.DefaultBasePath, client.DefaultSchemes)
	c := client.New(transport, nil)

	return url, shutdown, mockApp, c
}

// APIError returns model.Error with given msg.
func APIError(msg string) *models.Error {
	return &models.Error{
		Message: swag.String(msg),
	}
}

func errPayload(err interface{}) *models.Error {
	switch err := err.(type) {
	case *operations.VerificationEmailDefault:
		return err.Payload
	case *operations.VerificationUsernameDefault:
		return err.Payload
	case *operations.CreateUserDefault:
		return err.Payload
	case *operations.LoginDefault:
		return err.Payload
	case *operations.LogoutDefault:
		return err.Payload
	case *operations.GetUserDefault:
		return err.Payload
	case *operations.DeleteUserDefault:
		return err.Payload
	case *operations.UpdatePasswordDefault:
		return err.Payload
	case *operations.UpdateUsernameDefault:
		return err.Payload
	case *operations.UpdateEmailDefault:
		return err.Payload
	case *operations.GetUsersDefault:
		return err.Payload
	case *operations.CreateRecoveryCodeDefault:
		return err.Payload
	case *operations.RecoveryPasswordDefault:
		return err.Payload
	default:
		return nil
	}
}
