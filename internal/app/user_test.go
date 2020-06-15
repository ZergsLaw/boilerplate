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

	user := userGen(t)
	notExistEmail := notExistEmail

	mocks.userRepo.EXPECT().UserByEmail(ctx, notExistEmail).Return(nil, app.ErrNotFound)
	mocks.userRepo.EXPECT().UserByEmail(ctx, user.Email).Return(&user, nil)
	mocks.userRepo.EXPECT().UserByEmail(ctx, "").Return(nil, errAny)

	testCases := map[string]struct {
		email string
		want  error
	}{
		"success":   {notExistEmail, nil},
		"exist":     {user.Email, app.ErrEmailExist},
		"any error": {"", errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.VerificationEmail(ctx, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_VerificationUsername(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	user := userGen(t)
	mocks.userRepo.EXPECT().UserByUsername(ctx, notExistUsername).Return(nil, app.ErrNotFound)
	mocks.userRepo.EXPECT().UserByUsername(ctx, user.Name).Return(&user, nil)
	mocks.userRepo.EXPECT().UserByUsername(ctx, "").Return(nil, errAny)

	testCases := map[string]struct {
		username string
		want     error
	}{
		"success":   {notExistUsername, nil},
		"exist":     {user.Name, app.ErrUsernameExist},
		"any error": {"", errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.VerificationUsername(ctx, tc.username)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_Login(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	user := userGen(t)
	origin := newOrigin()

	const notValidPass = "notValidPass"
	const notValidTokenID app.TokenID = "notValidTokenID"
	// Couldn't come up with a proper name for the fucking var.
	notValidTokenExpireForGenerateNotValidTokenID := time.Second
	notValidTokenExpired := time.Second * 2

	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(user.Email)).Return(&user, nil).Times(4)
	mocks.password.EXPECT().Compare(user.PassHash, []byte(password)).Return(true).Times(3)
	mocks.auth.EXPECT().Token(app.TokenExpire).Return(token, tokenID, nil)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user.ID, tokenID, origin).Return(nil)

	mocks.auth.EXPECT().Token(notValidTokenExpireForGenerateNotValidTokenID).Return(token, notValidTokenID, nil)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user.ID, notValidTokenID, origin).Return(errAny)
	mocks.auth.EXPECT().Token(notValidTokenExpired).Return(app.AuthToken(""), app.TokenID(""), errAny)
	mocks.password.EXPECT().Compare(user.PassHash, []byte(notValidPass)).Return(false)
	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(notExistEmail)).Return(nil, app.ErrNotFound)

	testCases := map[string]struct {
		email       string
		password    string
		tokenExpire time.Duration
		want        *app.User
		wantToken   app.AuthToken
		wantErr     error
	}{
		"success":               {user.Email, password, app.TokenExpire, &user, token, nil},
		"err from save session": {user.Email, password, notValidTokenExpireForGenerateNotValidTokenID, nil, "", errAny},
		"err from gen token":    {user.Email, password, notValidTokenExpired, nil, "", errAny},
		"err from compare pass": {user.Email, notValidPass, app.TokenExpire, nil, "", app.ErrNotValidPassword},
		"user not found":        {notExistEmail, "", app.TokenExpire, nil, "", app.ErrNotFound},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			muTokenExpire.Lock()
			app.TokenExpire = tc.tokenExpire
			defer muTokenExpire.Unlock()

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

	const notValidEmail = "notValidEmail"
	const notCorrectPass = "notCorrectPass"
	user := userGen(t)
	origin := newOrigin()
	task := app.TaskNotification{
		Email: user.Email,
		Kind:  app.Welcome,
	}
	notValidTask := app.TaskNotification{
		Email: strings.ToLower(notValidEmail),
		Kind:  app.Welcome,
	}
	tokenExpire := 24 * 7 * time.Hour

	mocks.password.EXPECT().Hashing(password).Return([]byte(password), nil).Times(2)
	mocks.userRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    user.Email,
		Name:     user.Name,
		PassHash: []byte(password),
	}, task).Return(user.ID, nil)

	mocks.userRepo.EXPECT().UserByEmail(ctx, user.Email).Return(&user, nil)
	mocks.password.EXPECT().Compare(user.PassHash, []byte(password)).Return(true)
	mocks.auth.EXPECT().Token(tokenExpire).Return(token, tokenID, nil)
	mocks.sessionRepo.EXPECT().SaveSession(ctx, user.ID, tokenID, origin).Return(nil)
	mocks.userRepo.EXPECT().CreateUser(ctx, app.User{
		Email:    strings.ToLower(notValidEmail),
		Name:     user.Name,
		PassHash: []byte(password),
	}, notValidTask).Return(app.UserID(0), errAny)
	mocks.password.EXPECT().Hashing(notCorrectPass).Return(nil, errAny)

	testCases := map[string]struct {
		email     string
		username  string
		password  string
		want      *app.User
		wantToken app.AuthToken
		wantErr   error
	}{
		"success":         {user.Email, user.Name, password, &user, token, nil},
		"err create user": {notValidEmail, user.Name, password, nil, "", errAny},
		"err hashing":     {user.Email, user.Name, notCorrectPass, nil, "", errAny},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			muTokenExpire.Lock()
			app.TokenExpire = tokenExpire
			defer muTokenExpire.Unlock()

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

	user := userGen(t)
	mocks.userRepo.EXPECT().UpdateUsername(ctx, user.ID, notExistUsername).Return(nil)

	testCases := map[string]struct {
		username string
		want     error
	}{
		"success":         {notExistUsername, nil},
		"usernames equal": {user.Name, app.ErrUsernameNeedDifferentiate},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.UpdateUsername(ctx, app.AuthUser{User: user}, tc.username)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UpdateEmail(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	user := userGen(t)
	notExistEmail := notExistEmail
	task := app.TaskNotification{
		Email: strings.ToLower(notExistEmail),
		Kind:  app.ChangeEmail,
	}
	mocks.userRepo.EXPECT().UpdateEmail(ctx, user.ID, strings.ToLower(notExistEmail), task).Return(nil)

	testCases := map[string]struct {
		email string
		want  error
	}{
		"success":      {notExistEmail, nil},
		"emails equal": {user.Email, app.ErrEmailNeedDifferentiate},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.UpdateEmail(ctx, app.AuthUser{User: user}, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UpdatePassword(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	const notValidPass = "notValidPass"

	user := userGen(t)
	mocks.userRepo.EXPECT().UpdatePassword(ctx, user.ID, []byte(password)).Return(nil)
	mocks.password.EXPECT().Compare(user.PassHash, []byte(password)).Return(true).Times(2)
	mocks.password.EXPECT().Compare(user.PassHash, []byte(notValidPass)).Return(false).Times(1)
	mocks.password.EXPECT().Hashing(password).Return([]byte(password), nil)
	mocks.password.EXPECT().Hashing(notValidPass).Return(nil, errAny)

	testCases := map[string]struct {
		oldPass, newPass string
		want             error
	}{
		"success":                {password, password, nil},
		"err hashing":            {password, notValidPass, errAny},
		"err not valid password": {notValidPass, password, app.ErrNotValidPassword},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.UpdatePassword(ctx, app.AuthUser{User: user}, tc.oldPass, tc.newPass)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_CreateRecoveryCode(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	const codeLength = 6
	user := userGen(t)
	recoveryCode := recoveryCode
	notExistEmail := notExistEmail
	task := app.TaskNotification{
		Email: user.Email,
		Kind:  app.PassRecovery,
	}

	mocks.userRepo.EXPECT().UserByEmail(ctx, user.Email).Return(&user, nil)
	mocks.code.EXPECT().Generate(codeLength).Return(recoveryCode)
	mocks.codeRepo.EXPECT().SaveCode(ctx, user.Email, recoveryCode, task).Return(nil)
	mocks.userRepo.EXPECT().UserByEmail(ctx, strings.ToLower(notExistEmail)).Return(nil, app.ErrNotFound)

	testCases := map[string]struct {
		email string
		want  error
	}{
		"success":        {user.Email, nil},
		"user not found": {notExistEmail, app.ErrNotFound},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.CreateRecoveryCode(ctx, tc.email)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_RecoveryPassword(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	user := userGen(t)
	codeInfo := app.CodeInfo{
		Code:      recoveryCode,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}
	newPassword := "newPassword"
	notValidPass := "notValidPass"
	emailForExpiredCode := "expiredCode@test.test"
	emailForNotValidCode := "notValidCode@test.test"
	emailForNotExistCode := "notExistCode@test.test"

	mocks.userRepo.EXPECT().UserByEmail(ctx, user.Email).Return(&user, nil).Times(2)
	mocks.codeRepo.EXPECT().Code(ctx, user.Email).Return(&codeInfo, nil).Times(2)
	mocks.password.EXPECT().Hashing(newPassword).Return([]byte(newPassword), nil)
	mocks.userRepo.EXPECT().UpdatePassword(ctx, user.ID, []byte(newPassword)).Return(nil)

	mocks.password.EXPECT().Hashing(notValidPass).Return(nil, errAny)

	mocks.userRepo.EXPECT().UserByEmail(ctx, emailForExpiredCode).Return(&user, nil)
	mocks.codeRepo.EXPECT().Code(ctx, emailForExpiredCode).Return(&app.CodeInfo{
		Code:      codeInfo.Code,
		Email:     codeInfo.Email,
		CreatedAt: time.Time{},
	}, nil)

	mocks.userRepo.EXPECT().UserByEmail(ctx, emailForNotValidCode).Return(&user, nil)
	mocks.codeRepo.EXPECT().Code(ctx, emailForNotValidCode).Return(&app.CodeInfo{
		Code:      "any code",
		Email:     codeInfo.Email,
		CreatedAt: codeInfo.CreatedAt,
	}, nil)

	mocks.userRepo.EXPECT().UserByEmail(ctx, emailForNotExistCode).Return(&user, nil)
	mocks.codeRepo.EXPECT().Code(ctx, emailForNotExistCode).Return(nil, app.ErrNotFound)
	mocks.userRepo.EXPECT().UserByEmail(ctx, notExistEmail).Return(nil, app.ErrNotFound)

	testCases := map[string]struct {
		email   string
		newPass string
		want    error
	}{
		"success":           {user.Email, newPassword, nil},
		"err from hashing":  {user.Email, notValidPass, errAny},
		"expired":           {emailForExpiredCode, "", app.ErrCodeExpired},
		"not valid":         {emailForNotValidCode, "", app.ErrNotValidCode},
		"err from get code": {emailForNotExistCode, "", app.ErrNotFound},
		"err from get user": {notExistEmail, "", app.ErrNotFound},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			err := application.RecoveryPassword(ctx, tc.email, recoveryCode, tc.newPass)
			assert.Equal(t, tc.want, err)
		})
	}
}

func TestApp_UserByAuthToken(t *testing.T) {
	t.Parallel()

	application, mocks, shutdown := initTest(t)
	defer shutdown()

	const expiredToken app.AuthToken = "notValidToken"

	user := userGen(t)
	session := sessionGen(t)
	auth := app.AuthUser{
		User:    user,
		Session: session,
	}

	mocks.auth.EXPECT().Parse(token).Return(tokenID, nil).Times(3)
	mocks.auth.EXPECT().Parse(expiredToken).Return(app.TokenID(""), app.ErrExpiredToken)
	mocks.sessionRepo.EXPECT().UserByTokenID(ctx, tokenID).Return(&user, nil).Times(2)
	mocks.sessionRepo.EXPECT().UserByTokenID(ctx, tokenID).Return(nil, app.ErrNotFound)
	mocks.sessionRepo.EXPECT().SessionByTokenID(ctx, tokenID).Return(&session, nil)
	mocks.sessionRepo.EXPECT().SessionByTokenID(ctx, tokenID).Return(nil, errAny)

	testCases := []struct {
		name    string
		token   app.AuthToken
		want    *app.AuthUser
		wantErr error
	}{
		{"success", token, &auth, nil},
		{"invalid token", "", nil, app.ErrInvalidToken},
		{"err session by auth", token, nil, errAny},
		{"not found user by auth", token, nil, app.ErrNotFound},
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
