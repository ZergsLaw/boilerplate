package app_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestApp_VerificationEmail(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UserByEmail(ctx, notExistEmail).Return(nil, app.ErrNotFound)
	mocks.userRepo.EXPECT().UserByEmail(ctx, email1).Return(&user1, nil)
	mocks.userRepo.EXPECT().UserByEmail(ctx, "").Return(nil, errAny)

	testCases := []struct {
		name  string
		email string
		want  error
	}{
		{"success", notExistEmail, nil},
		{"exist", email1, app.ErrEmailExist},
		{"any error", "", errAny},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.VerificationEmail(ctx, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_VerificationUsername(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UserByUsername(ctx, notExistUsername).Return(nil, app.ErrNotFound)
	mocks.userRepo.EXPECT().UserByUsername(ctx, username).Return(&user1, nil)
	mocks.userRepo.EXPECT().UserByUsername(ctx, "").Return(nil, errAny)

	testCases := []struct {
		name     string
		username string
		want     error
	}{
		{"success", notExistUsername, nil},
		{"exist", username, app.ErrUsernameExist},
		{"any error", "", errAny},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.VerificationUsername(ctx, tc.username)
			assert.Equal(t, tc.want, err)
		})
	}
}

const tokenExpire = 24 * 7 * time.Hour

func TestApp_Login(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(email1)).Return(&user1, nil).Times(4)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(nil)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(errAny)
	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(notExistEmail)).Return(nil, app.ErrNotFound)
	mocks.password.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true).Times(3)
	mocks.password.EXPECT().Compare(user1.PassHash, []byte(password2)).Return(false)
	mocks.auth.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)
	mocks.auth.EXPECT().Token(tokenExpire).Return(app.AuthToken(""), app.TokenID(""), errAny)
	mocks.auth.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)

	testCases := []struct {
		name      string
		email     string
		password  string
		want      *app.User
		wantToken app.AuthToken
		wantErr   error
	}{
		{"success", email1, password1, &user1, token1, nil},
		{"user not found", notExistEmail, "", nil, "", app.ErrNotFound},
		{"not correct password", email1, password2, nil, "", app.ErrNotValidPassword},
		{"error generate token1", email1, password1, nil, "", errAny},
		{"not save session", email1, password1, nil, "", errAny},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			user, token, err := application.Login(ctx, tc.email, tc.password, origin)
			if tc.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, user)
				assert.Equal(t, tc.wantToken, token)
			} else {
				assert.Nil(t, user)
				assert.Equal(t, tc.wantToken, token)
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}

func TestApp_CreateUser(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.password.EXPECT().Hashing(password1).Return([]byte(password1), nil).Times(2)
	mocks.userRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    email1,
		Username: username,
		PassHash: []byte(password1),
	}).Return(user1.ID, nil)
	mocks.userRepo.EXPECT().UserByEmail(ctx, email1).Return(&user1, nil)
	mocks.password.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true)
	mocks.auth.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(nil)

	mocks.userRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    email1,
		Username: username,
		PassHash: []byte(password1),
	}).Return(app.UserID(0), errAny)

	mocks.password.EXPECT().Hashing(password1).Return(nil, errAny)

	testCases := []struct {
		name      string
		email     string
		username  string
		password  string
		want      *app.User
		wantToken app.AuthToken
		wantErr   error
	}{
		{"success", email1, username, password1, &user1, token1, nil},
		{"err create user", email1, username, password1, nil, "", errAny},
		{"err hashing", email1, username, password1, nil, "", errAny},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			user, token, err := application.CreateUser(ctx, tc.email, tc.username, tc.password, origin)
			if tc.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, user)
				assert.Equal(t, tc.wantToken, token)
			} else {
				assert.Nil(t, user)
				assert.Equal(t, tc.wantToken, token)
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}

func TestApp_UpdateUsername(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UpdateUsername(ctx, user1.ID, notExistUsername).Return(nil)

	testCases := []struct {
		name     string
		username string
		want     error
	}{
		{"success", notExistUsername, nil},
		{"usernames equal", username, app.ErrUsernameNeedDifferentiate},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.UpdateUsername(ctx, app.AuthUser{User: user1}, tc.username)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UpdateEmail(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UpdateEmail(ctx, user1.ID, strings.ToLower(notExistEmail)).Return(nil)

	testCases := []struct {
		name  string
		email string
		want  error
	}{
		{"success", notExistEmail, nil},
		{"emails equal", email1, app.ErrEmailNeedDifferentiate},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.UpdateEmail(ctx, app.AuthUser{User: user1}, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UpdatePassword(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.userRepo.EXPECT().UpdatePassword(ctx, user1.ID, []byte(password2)).Return(nil)
	mocks.password.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true).Times(2)
	mocks.password.EXPECT().Compare(user1.PassHash, []byte(password2)).Return(false).Times(1)
	mocks.password.EXPECT().Hashing(password2).Return([]byte(password2), nil)
	mocks.password.EXPECT().Hashing(password2).Return(nil, errAny)

	testCases := []struct {
		name             string
		oldPass, newPass string
		want             error
	}{
		{"success", password1, password2, nil},
		{"err hashing", password1, password2, errAny},
		{"err not valid password", password2, password2, app.ErrNotValidPassword},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.UpdatePassword(ctx, app.AuthUser{User: user1}, tc.oldPass, tc.newPass)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_CreateRecoveryCode(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	const codeLength = 6
	mocks.userRepo.EXPECT().UserByEmail(ctx, email1).Return(&user1, nil)
	mocks.code.EXPECT().Generate(codeLength).Return(recoveryCode)
	mocks.codeRepo.EXPECT().SaveCode(ctx, user1.ID, recoveryCode)
	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(notExistEmail)).Return(nil, app.ErrNotFound)

	testCases := []struct {
		name  string
		email string
		want  error
	}{
		{"success", email1, nil},
		{"user not found", notExistEmail, app.ErrNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.CreateRecoveryCode(ctx, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_RecoveryPassword(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	mocks.codeRepo.EXPECT().UserIDByCode(ctx, recoveryCode).Return(user1.ID, time.Now(), nil).Times(2)
	mocks.password.EXPECT().Hashing(password2).Return([]byte(password2), nil)
	mocks.userRepo.EXPECT().UpdatePassword(ctx, user1.ID, []byte(password2)).Return(nil)
	mocks.password.EXPECT().Hashing(password2).Return(nil, errAny)
	mocks.codeRepo.EXPECT().UserIDByCode(ctx, recoveryCode).Return(user1.ID, time.Time{}, nil)
	mocks.codeRepo.EXPECT().UserIDByCode(ctx, recoveryCode).Return(app.UserID(0), time.Time{}, app.ErrNotFound)

	testCases := []struct {
		name string
		want error
	}{
		{"success", nil},
		{"hashing error", errAny},
		{"recovery recoverycode is expired", app.ErrCodeExpired},
		{"not found email by recoverycode", app.ErrNotFound},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := application.RecoveryPassword(ctx, recoveryCode, password2)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UserByAuthToken(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	const expiredToken app.AuthToken = "notValidToken"

	mocks.auth.EXPECT().Parse(token1).Return(tokenID1, nil).Times(3)
	mocks.auth.EXPECT().Parse(expiredToken).Return(app.TokenID(""), app.ErrExpiredToken)
	mocks.sessionRepo.EXPECT().UserByTokenID(ctx, tokenID1).Return(&user1, nil).Times(2)
	mocks.sessionRepo.EXPECT().UserByTokenID(ctx, tokenID1).Return(nil, app.ErrNotFound)
	mocks.sessionRepo.EXPECT().SessionByTokenID(ctx, tokenID1).Return(&session1, nil)
	mocks.sessionRepo.EXPECT().SessionByTokenID(ctx, tokenID1).Return(nil, errAny)

	testCases := []struct {
		name    string
		token   app.AuthToken
		want    *app.AuthUser
		wantErr error
	}{
		{"success", token1, &authUser, nil},
		{"invalid token", "", nil, app.ErrInvalidToken},
		{"err session by auth", token1, nil, errAny},
		{"not found user by auth", token1, nil, app.ErrNotFound},
		{"not valid auth", expiredToken, nil, app.ErrExpiredToken},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			user, err := application.UserByAuthToken(ctx, tc.token)
			if tc.wantErr == nil {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, user)
			} else {
				assert.Nil(t, user)
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}
