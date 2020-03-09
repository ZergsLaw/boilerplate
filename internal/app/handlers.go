package app

import (
	"context"
	"errors"
	"strings"
	"time"
)

func (a *app) VerificationEmail(ctx context.Context, email string) error {
	_, err := a.repo.UserByEmail(ctx, email)
	switch {
	case errors.Is(err, ErrNotFound):
		return nil
	case err == nil:
		return ErrEmailExist
	default:
		return err
	}
}

func (a *app) VerificationUsername(ctx context.Context, username string) error {
	_, err := a.repo.UserByUsername(ctx, username)
	switch {
	case errors.Is(err, ErrNotFound):
		return nil
	case err == nil:
		return ErrUsernameExist
	default:
		return err
	}
}

const (
	tokenExpire = 24 * 7 * time.Hour
)

func (a *app) Login(ctx context.Context, email, password string, origin Origin) (*User, AuthToken, error) {
	email = strings.ToLower(email)

	user, err := a.repo.UserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if !a.password.Compare(user.PassHash, []byte(password)) {
		return nil, "", ErrNotValidPassword
	}

	token, tokenID, err := a.auth.Token(tokenExpire)
	if err != nil {
		return nil, "", err
	}

	err = a.repo.SaveSession(ctx, user.ID, tokenID, origin)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (a *app) Logout(ctx context.Context, authUser AuthUser) error {
	return a.repo.DeleteSession(ctx, authUser.Session.TokenID)
}

func (a *app) CreateUser(ctx context.Context, email, username, password string, origin Origin) (*User, AuthToken, error) {
	passHash, err := a.password.Hashing(password)
	if err != nil {
		return nil, "", err
	}
	email = strings.ToLower(email)

	_, err = a.repo.CreateUser(ctx, User{
		Email:    email,
		Username: username,
		PassHash: passHash,
	})
	if err != nil {
		return nil, "", err
	}

	return a.Login(ctx, email, password, origin)
}

func (a *app) User(ctx context.Context, _ AuthUser, userID UserID) (*User, error) {
	return a.repo.UserByID(ctx, userID)
}

func (a *app) DeleteUser(ctx context.Context, authUser AuthUser) error {
	return a.repo.DeleteUser(ctx, authUser.ID)
}

func (a *app) UpdateUsername(ctx context.Context, authUser AuthUser, username string) error {
	if authUser.Username == username {
		return ErrUsernameNeedDifferentiate
	}

	return a.repo.UpdateUsername(ctx, authUser.ID, username)
}

func (a *app) UpdateEmail(ctx context.Context, authUser AuthUser, email string) error {
	email = strings.ToLower(email)
	if authUser.Email == email {
		return ErrEmailNeedDifferentiate
	}

	return a.repo.UpdateEmail(ctx, authUser.ID, email)
}

func (a *app) UpdatePassword(ctx context.Context, authUser AuthUser, oldPass, newPass string) error {
	if !a.password.Compare(authUser.PassHash, []byte(oldPass)) {
		return ErrNotValidPassword
	}

	passHash, err := a.password.Hashing(newPass)
	if err != nil {
		return err
	}

	return a.repo.UpdatePassword(ctx, authUser.ID, passHash)
}

func (a *app) ListUserByUsername(ctx context.Context, _ AuthUser, username string, page Page) ([]User, int, error) {
	return a.repo.ListUserByUsername(ctx, username, page)
}

func (a *app) UserByAuthToken(ctx context.Context, token AuthToken) (*AuthUser, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	tokenID, err := a.auth.Parse(token)
	if err != nil {
		return nil, err
	}

	user, err := a.repo.UserByTokenID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	session, err := a.repo.SessionByTokenID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	return &AuthUser{User: *user, Session: *session}, nil
}
