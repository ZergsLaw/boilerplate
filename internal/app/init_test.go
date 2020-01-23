package app_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/mock"
)

var (
	ctx    = context.Background()
	errAny = errors.New("any error")

	notExistEmail = "notExist@email.com"
	email1        = "exist@email1.com"

	notExistUsername = "notExistUsername"
	username         = "username"

	password1 = "password1"
	password2 = "password2"

	token1 app.AuthToken = "token1"

	tokenID1 app.TokenID = "tokenID1"

	session1 = app.Session{
		Origin:  origin,
		ID:      1,
		TokenID: tokenID1,
	}

	origin = app.Origin{
		IP:        net.ParseIP("192.100.10.4"),
		UserAgent: "UserAgent",
	}

	user1 = app.User{
		ID:        1,
		Email:     email1,
		Username:  username,
		PassHash:  []byte(password1),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	authUser = app.AuthUser{
		User:    user1,
		Session: session1,
	}
)

func initTest(t testing.TB) (app.App, *mock.Repo, *mock.Password, *mock.Auth, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockRepo := mock.NewRepo(ctrl)
	mockPass := mock.NewPassword(ctrl)
	mockToken := mock.NewAuth(ctrl)

	return app.New(mockRepo, mockPass, mockToken),
		mockRepo, mockPass, mockToken, ctrl.Finish
}
