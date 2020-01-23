package grpc

import (
	"context"
	"errors"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) GetUserByAuthToken(ctx context.Context, auth *AuthInfo) (*User, error) {
	logger := log.FromContext(ctx)
	logger.Infof("get user profile by: %s", auth.Token)

	info, err := s.app.UserByAuthToken(ctx, app.AuthToken(auth.Token))
	if err != nil {
		return nil, apiError(err)
	}

	return apiUser(&info.User), nil
}

func (s *service) GetUserByID(ctx context.Context, userID *UserID) (*User, error) {
	logger := log.FromContext(ctx)
	logger.Infof("get user profile by: %d", userID.Id)

	user, err := s.app.User(ctx, app.UserID(userID.Id))
	if err != nil {
		return nil, apiError(err)
	}

	return apiUser(user), nil
}

func apiUser(user *app.User) *User {
	return &User{
		Id:       int32(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}
}

func apiError(err error) error {
	if err == nil {
		return nil
	}

	code := codes.Internal
	switch {
	case errors.Is(err, app.ErrNotFound):
		code = codes.NotFound
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	}

	return status.Error(code, err.Error())
}
