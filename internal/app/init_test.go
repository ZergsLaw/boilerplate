package app_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/mock"
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

	taskNotification = app.TaskNotification{
		ID:     1,
		UserID: user1.ID,
		Kind:   app.Welcome,
	}

	recoveryCode = "123456"
)

type Mocks struct {
	userRepo     *mock.MockUserRepo
	sessionRepo  *mock.MockSessionRepo
	codeRepo     *mock.MockCodeRepo
	code         *mock.MockCode
	password     *mock.MockPassword
	auth         *mock.MockAuth
	wal          *mock.MockWAL
	notification *mock.MockNotification
}

func initTest(t *testing.T) (*app.Application, *Mocks, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)

	mockUserRepo := mock.NewMockUserRepo(ctrl)
	mockSessionRepo := mock.NewMockSessionRepo(ctrl)
	mockCodeRepo := mock.NewMockCodeRepo(ctrl)
	mockCode := mock.NewMockCode(ctrl)
	mockPass := mock.NewMockPassword(ctrl)
	mockToken := mock.NewMockAuth(ctrl)
	mockWal := mock.NewMockWAL(ctrl)
	mockNotification := mock.NewMockNotification(ctrl)

	appl := app.New(app.Config{
		UserRepo:     mockUserRepo,
		SessionRepo:  mockSessionRepo,
		CodeRepo:     mockCodeRepo,
		Password:     mockPass,
		Auth:         mockToken,
		Wal:          mockWal,
		Notification: mockNotification,
		Code:         mockCode,
	})

	mocks := &Mocks{
		userRepo:     mockUserRepo,
		sessionRepo:  mockSessionRepo,
		codeRepo:     mockCodeRepo,
		code:         mockCode,
		password:     mockPass,
		auth:         mockToken,
		wal:          mockWal,
		notification: mockNotification,
	}

	return appl, mocks, ctrl.Finish
}
