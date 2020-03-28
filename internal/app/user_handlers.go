package app

import (
	"context"
	"errors"
	"strings"
	"time"
)

// VerificationEmail for implemented UserApp.
func (a *Application) VerificationEmail(ctx context.Context, email string) error {
	_, err := a.userRepo.UserByEmail(ctx, email)
	switch {
	case errors.Is(err, ErrNotFound):
		return nil
	case err == nil:
		return ErrEmailExist
	default:
		return err
	}
}

// VerificationUsername for implemented UserApp.
func (a *Application) VerificationUsername(ctx context.Context, username string) error {
	_, err := a.userRepo.UserByUsername(ctx, username)
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

// Login for implemented UserApp.
func (a *Application) Login(ctx context.Context, email, password string, origin Origin) (*User, AuthToken, error) {
	email = strings.ToLower(email)

	user, err := a.userRepo.UserByEmail(ctx, email)
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

	err = a.sessionRepo.SaveSession(ctx, user.ID, tokenID, origin)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Logout for implemented UserApp.
func (a *Application) Logout(ctx context.Context, authUser AuthUser) error {
	return a.sessionRepo.DeleteSession(ctx, authUser.Session.TokenID)
}

// CreateUser for implemented UserApp.
func (a *Application) CreateUser(ctx context.Context, email, username, password string, origin Origin) (*User, AuthToken, error) {
	passHash, err := a.password.Hashing(password)
	if err != nil {
		return nil, "", err
	}
	email = strings.ToLower(email)

	_, err = a.userRepo.CreateUser(ctx, User{
		Email:    email,
		Username: username,
		PassHash: passHash,
	})
	if err != nil {
		return nil, "", err
	}

	return a.Login(ctx, email, password, origin)
}

// User for implemented UserApp.
func (a *Application) User(ctx context.Context, _ AuthUser, userID UserID) (*User, error) {
	return a.userRepo.UserByID(ctx, userID)
}

// DeleteUser for implemented UserApp.
func (a *Application) DeleteUser(ctx context.Context, authUser AuthUser) error {
	return a.userRepo.DeleteUser(ctx, authUser.ID)
}

// UpdateUsername for implemented UserApp.
func (a *Application) UpdateUsername(ctx context.Context, authUser AuthUser, username string) error {
	if authUser.Username == username {
		return ErrUsernameNeedDifferentiate
	}

	return a.userRepo.UpdateUsername(ctx, authUser.ID, username)
}

// UpdateEmail for implemented UserApp.
func (a *Application) UpdateEmail(ctx context.Context, authUser AuthUser, email string) error {
	email = strings.ToLower(email)
	if authUser.Email == email {
		return ErrEmailNeedDifferentiate
	}

	return a.userRepo.UpdateEmail(ctx, authUser.ID, email)
}

// UpdatePassword for implemented UserApp.
func (a *Application) UpdatePassword(ctx context.Context, authUser AuthUser, oldPass, newPass string) error {
	if !a.password.Compare(authUser.PassHash, []byte(oldPass)) {
		return ErrNotValidPassword
	}

	passHash, err := a.password.Hashing(newPass)
	if err != nil {
		return err
	}

	return a.userRepo.UpdatePassword(ctx, authUser.ID, passHash)
}

// ListUserByUsername for implemented UserApp.
func (a *Application) ListUserByUsername(ctx context.Context, _ AuthUser, username string, page Page) ([]User, int, error) {
	return a.userRepo.ListUserByUsername(ctx, username, page)
}

// CreateRecoveryCode for implemented UserApp.
func (a *Application) CreateRecoveryCode(ctx context.Context, email string) error {
	const codeLength = 6
	email = strings.ToLower(email)

	user, err := a.userRepo.UserByEmail(ctx, email)
	if err != nil {
		return err
	}

	code := a.code.Generate(codeLength)

	return a.codeRepo.SaveCode(ctx, user.ID, code)
}

// RecoveryPassword for implemented UserApp.
func (a *Application) RecoveryPassword(ctx context.Context, code, newPassword string) error {
	userID, createdAt, err := a.codeRepo.UserID(ctx, code)
	if err != nil {
		return err
	}

	const recoveryCodeLifetime = time.Hour * 24
	if time.Since(createdAt) > recoveryCodeLifetime {
		return ErrCodeExpired
	}

	passHash, err := a.password.Hashing(newPassword)
	if err != nil {
		return err
	}

	return a.userRepo.UpdatePassword(ctx, userID, passHash)
}

// UserByAuthToken for implemented UserApp.
func (a *Application) UserByAuthToken(ctx context.Context, token AuthToken) (*AuthUser, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	tokenID, err := a.auth.Parse(token)
	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.UserByTokenID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	session, err := a.sessionRepo.SessionByTokenID(ctx, tokenID)
	if err != nil {
		return nil, err
	}

	return &AuthUser{User: *user, Session: *session}, nil
}
