package app_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/users/internal/app"
)

func TestApp_VerificationEmail(t *testing.T) {
	t.Parallel()

	application, mockRepo, _, _, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UserByEmail(ctx, notExistEmail).Return(nil, app.ErrNotFound)
	mockRepo.EXPECT().UserByEmail(ctx, email1).Return(&user1, nil)
	mockRepo.EXPECT().UserByEmail(ctx, "").Return(nil, errAny)

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

	application, mockRepo, _, _, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UserByUsername(ctx, notExistUsername).Return(nil, app.ErrNotFound)
	mockRepo.EXPECT().UserByUsername(ctx, username).Return(&user1, nil)
	mockRepo.EXPECT().UserByUsername(ctx, "").Return(nil, errAny)

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

	application, mockRepo, mockPassword, mockToken, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UserByEmail(ctx, strings.ToLower(email1)).Return(&user1, nil).Times(4)
	mockRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(nil)
	mockRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(errAny)
	mockRepo.EXPECT().UserByEmail(ctx, strings.ToLower(notExistEmail)).Return(nil, app.ErrNotFound)
	mockPassword.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true).Times(3)
	mockPassword.EXPECT().Compare(user1.PassHash, []byte(password2)).Return(false)
	mockToken.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)
	mockToken.EXPECT().Token(tokenExpire).Return(app.AuthToken(""), app.TokenID(""), errAny)
	mockToken.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)

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

	application, mockRepo, mockPassword, mockToken, shutdown := initTest(t)
	defer shutdown()

	mockPassword.EXPECT().Hashing(password1).Return([]byte(password1), nil).Times(2)
	mockRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    email1,
		Username: username,
		PassHash: []byte(password1),
	}).Return(user1.ID, nil)
	mockRepo.EXPECT().UserByEmail(ctx, email1).Return(&user1, nil)
	mockPassword.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true)
	mockToken.EXPECT().Token(tokenExpire).Return(token1, tokenID1, nil)
	mockRepo.EXPECT().SaveSession(ctx, user1.ID, tokenID1, origin).Return(nil)

	mockRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    email1,
		Username: username,
		PassHash: []byte(password1),
	}).Return(app.UserID(0), errAny)

	mockPassword.EXPECT().Hashing(password1).Return(nil, errAny)

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

	application, mockRepo, _, _, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UpdateUsername(ctx, user1.ID, notExistUsername).Return(nil)

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

	application, mockRepo, _, _, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UpdateEmail(ctx, user1.ID, strings.ToLower(notExistEmail)).Return(nil)

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

	application, mockRepo, mockPassword, _, shutdown := initTest(t)
	defer shutdown()

	mockRepo.EXPECT().UpdatePassword(ctx, user1.ID, []byte(password2)).Return(nil)
	mockPassword.EXPECT().Compare(user1.PassHash, []byte(password1)).Return(true).Times(2)
	mockPassword.EXPECT().Compare(user1.PassHash, []byte(password2)).Return(false).Times(1)
	mockPassword.EXPECT().Hashing(password2).Return([]byte(password2), nil)
	mockPassword.EXPECT().Hashing(password2).Return(nil, errAny)

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

func TestApp_UserByAuthToken(t *testing.T) {
	t.Parallel()

	application, mockRepo, _, mockToken, shutdown := initTest(t)
	defer shutdown()

	const expiredToken app.AuthToken = "notValidToken"

	mockToken.EXPECT().Parse(token1).Return(tokenID1, nil).Times(3)
	mockToken.EXPECT().Parse(expiredToken).Return(app.TokenID(""), app.ErrExpiredToken)
	mockRepo.EXPECT().UserByTokenID(ctx, tokenID1).Return(&user1, nil).Times(2)
	mockRepo.EXPECT().UserByTokenID(ctx, tokenID1).Return(nil, app.ErrNotFound)
	mockRepo.EXPECT().SessionByTokenID(ctx, tokenID1).Return(&session1, nil)
	mockRepo.EXPECT().SessionByTokenID(ctx, tokenID1).Return(nil, errAny)

	testCases := []struct {
		name    string
		token   app.AuthToken
		want    *app.AuthUser
		wantErr error
	}{
		{"success", token1, &authUser, nil},
		{"success", "", nil, app.ErrInvalidToken},
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
