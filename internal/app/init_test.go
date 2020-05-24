package app_test

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zergslaw/boilerplate/internal/app"
	"github.com/zergslaw/boilerplate/internal/mock"
)

const (
	username  = "username"
	userEmail = "email@email.com"

	notExistEmail    = "notExist@email.com"
	notExistUsername = "notExistUsername"

	password = "password"

	token   app.AuthToken = "token"
	tokenID app.TokenID   = "tokenID"

	recoveryCode = "123456"

	ip        = "192.100.10.4"
	userAgent = "UserAgent"
)

var (
	ctx        = context.Background()
	errAny     = errors.New("any error")
	userGen    = userGenerator()
	sessionGen = sessionGenerator()

	// For def app.TokenExpire in test.
	muTokenExpire = sync.Mutex{}
)

func userGenerator() func(t *testing.T) app.User {
	x := app.UserID(0)
	mu := sync.Mutex{}

	return func(t *testing.T) app.User {
		t.Helper()

		mu.Lock()
		defer mu.Unlock()
		x++

		xStr := strconv.Itoa(int(x))
		return app.User{
			ID:        x,
			Email:     userEmail + xStr,
			Name:      username + xStr,
			PassHash:  []byte(password + xStr),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
}

func sessionGenerator() func(t *testing.T) app.Session {
	x := app.SessionID(0)
	mu := sync.Mutex{}

	return func(t *testing.T) app.Session {
		t.Helper()

		mu.Lock()
		defer mu.Unlock()
		x++

		xStr := strconv.Itoa(int(x))
		return app.Session{
			Origin:  newOrigin(),
			ID:      x,
			TokenID: tokenID + app.TokenID(xStr),
		}
	}
}

func newSession() app.Session {
	return app.Session{
		Origin:  newOrigin(),
		ID:      1,
		TokenID: tokenID,
	}
}

func newOrigin() app.Origin {
	return app.Origin{
		IP:        net.ParseIP(ip),
		UserAgent: userAgent,
	}
}

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
