package rpc

import (
	"context"
	"errors"

	"github.com/zergslaw/users/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) GetUserByAuthToken(ctx context.Context, in *AuthInfo) (*User, error) {
	info, err := s.app.UserByAuthToken(ctx, app.AuthToken(in.Token))
	if err != nil {
		return nil, apiError(err)
	}

	return apiUser(&info.User), nil
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
