// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/user.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	app "github.com/zergslaw/boilerplate/internal/app"
	reflect "reflect"
	time "time"
)

// MockUserApp is a mock of UserApp interface
type MockUserApp struct {
	ctrl     *gomock.Controller
	recorder *MockUserAppMockRecorder
}

// MockUserAppMockRecorder is the mock recorder for MockUserApp
type MockUserAppMockRecorder struct {
	mock *MockUserApp
}

// NewMockUserApp creates a new mock instance
func NewMockUserApp(ctrl *gomock.Controller) *MockUserApp {
	mock := &MockUserApp{ctrl: ctrl}
	mock.recorder = &MockUserAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserApp) EXPECT() *MockUserAppMockRecorder {
	return m.recorder
}

// VerificationEmail mocks base method
func (m *MockUserApp) VerificationEmail(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerificationEmail", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerificationEmail indicates an expected call of VerificationEmail
func (mr *MockUserAppMockRecorder) VerificationEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerificationEmail", reflect.TypeOf((*MockUserApp)(nil).VerificationEmail), ctx, email)
}

// VerificationUsername mocks base method
func (m *MockUserApp) VerificationUsername(ctx context.Context, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerificationUsername", ctx, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerificationUsername indicates an expected call of VerificationUsername
func (mr *MockUserAppMockRecorder) VerificationUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerificationUsername", reflect.TypeOf((*MockUserApp)(nil).VerificationUsername), ctx, username)
}

// Login mocks base method
func (m *MockUserApp) Login(ctx context.Context, email, password string, origin app.Origin) (*app.User, app.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password, origin)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(app.AuthToken)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Login indicates an expected call of Login
func (mr *MockUserAppMockRecorder) Login(ctx, email, password, origin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserApp)(nil).Login), ctx, email, password, origin)
}

// Logout mocks base method
func (m *MockUserApp) Logout(arg0 context.Context, arg1 app.AuthUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout
func (mr *MockUserAppMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockUserApp)(nil).Logout), arg0, arg1)
}

// CreateUser mocks base method
func (m *MockUserApp) CreateUser(ctx context.Context, email, username, password string, origin app.Origin) (*app.User, app.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, email, username, password, origin)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(app.AuthToken)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserAppMockRecorder) CreateUser(ctx, email, username, password, origin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserApp)(nil).CreateUser), ctx, email, username, password, origin)
}

// DeleteUser mocks base method
func (m *MockUserApp) DeleteUser(arg0 context.Context, arg1 app.AuthUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUserAppMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserApp)(nil).DeleteUser), arg0, arg1)
}

// User mocks base method
func (m *MockUserApp) User(arg0 context.Context, arg1 app.AuthUser, arg2 app.UserID) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", arg0, arg1, arg2)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User
func (mr *MockUserAppMockRecorder) User(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockUserApp)(nil).User), arg0, arg1, arg2)
}

// UserByAuthToken mocks base method
func (m *MockUserApp) UserByAuthToken(ctx context.Context, token app.AuthToken) (*app.AuthUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByAuthToken", ctx, token)
	ret0, _ := ret[0].(*app.AuthUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByAuthToken indicates an expected call of UserByAuthToken
func (mr *MockUserAppMockRecorder) UserByAuthToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByAuthToken", reflect.TypeOf((*MockUserApp)(nil).UserByAuthToken), ctx, token)
}

// UpdateUsername mocks base method
func (m *MockUserApp) UpdateUsername(arg0 context.Context, arg1 app.AuthUser, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUsername", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUsername indicates an expected call of UpdateUsername
func (mr *MockUserAppMockRecorder) UpdateUsername(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUsername", reflect.TypeOf((*MockUserApp)(nil).UpdateUsername), arg0, arg1, arg2)
}

// UpdateEmail mocks base method
func (m *MockUserApp) UpdateEmail(arg0 context.Context, arg1 app.AuthUser, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail
func (mr *MockUserAppMockRecorder) UpdateEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockUserApp)(nil).UpdateEmail), arg0, arg1, arg2)
}

// UpdatePassword mocks base method
func (m *MockUserApp) UpdatePassword(ctx context.Context, authUser app.AuthUser, oldPass, newPass string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", ctx, authUser, oldPass, newPass)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword
func (mr *MockUserAppMockRecorder) UpdatePassword(ctx, authUser, oldPass, newPass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserApp)(nil).UpdatePassword), ctx, authUser, oldPass, newPass)
}

// ListUserByUsername mocks base method
func (m *MockUserApp) ListUserByUsername(arg0 context.Context, arg1 app.AuthUser, arg2 string, arg3 app.Page) ([]app.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserByUsername", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]app.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUserByUsername indicates an expected call of ListUserByUsername
func (mr *MockUserAppMockRecorder) ListUserByUsername(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserByUsername", reflect.TypeOf((*MockUserApp)(nil).ListUserByUsername), arg0, arg1, arg2, arg3)
}

// CreateRecoveryCode mocks base method
func (m *MockUserApp) CreateRecoveryCode(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecoveryCode", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRecoveryCode indicates an expected call of CreateRecoveryCode
func (mr *MockUserAppMockRecorder) CreateRecoveryCode(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecoveryCode", reflect.TypeOf((*MockUserApp)(nil).CreateRecoveryCode), ctx, email)
}

// RecoveryPassword mocks base method
func (m *MockUserApp) RecoveryPassword(ctx context.Context, code, newPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecoveryPassword", ctx, code, newPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecoveryPassword indicates an expected call of RecoveryPassword
func (mr *MockUserAppMockRecorder) RecoveryPassword(ctx, code, newPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoveryPassword", reflect.TypeOf((*MockUserApp)(nil).RecoveryPassword), ctx, code, newPassword)
}

// MockUserRepo is a mock of UserRepo interface
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CreateUser mocks base method
func (m *MockUserRepo) CreateUser(arg0 context.Context, arg1 app.User) (app.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(app.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockUserRepoMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), arg0, arg1)
}

// DeleteUser mocks base method
func (m *MockUserRepo) DeleteUser(arg0 context.Context, arg1 app.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUserRepoMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepo)(nil).DeleteUser), arg0, arg1)
}

// UpdateUsername mocks base method
func (m *MockUserRepo) UpdateUsername(arg0 context.Context, arg1 app.UserID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUsername", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUsername indicates an expected call of UpdateUsername
func (mr *MockUserRepoMockRecorder) UpdateUsername(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUsername", reflect.TypeOf((*MockUserRepo)(nil).UpdateUsername), arg0, arg1, arg2)
}

// UpdateEmail mocks base method
func (m *MockUserRepo) UpdateEmail(arg0 context.Context, arg1 app.UserID, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail
func (mr *MockUserRepoMockRecorder) UpdateEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockUserRepo)(nil).UpdateEmail), arg0, arg1, arg2)
}

// UpdatePassword mocks base method
func (m *MockUserRepo) UpdatePassword(arg0 context.Context, arg1 app.UserID, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword
func (mr *MockUserRepoMockRecorder) UpdatePassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockUserRepo)(nil).UpdatePassword), arg0, arg1, arg2)
}

// UserByID mocks base method
func (m *MockUserRepo) UserByID(arg0 context.Context, arg1 app.UserID) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByID", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByID indicates an expected call of UserByID
func (mr *MockUserRepoMockRecorder) UserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByID", reflect.TypeOf((*MockUserRepo)(nil).UserByID), arg0, arg1)
}

// UserByEmail mocks base method
func (m *MockUserRepo) UserByEmail(arg0 context.Context, arg1 string) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByEmail indicates an expected call of UserByEmail
func (mr *MockUserRepoMockRecorder) UserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByEmail", reflect.TypeOf((*MockUserRepo)(nil).UserByEmail), arg0, arg1)
}

// UserByUsername mocks base method
func (m *MockUserRepo) UserByUsername(arg0 context.Context, arg1 string) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByUsername", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByUsername indicates an expected call of UserByUsername
func (mr *MockUserRepoMockRecorder) UserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByUsername", reflect.TypeOf((*MockUserRepo)(nil).UserByUsername), arg0, arg1)
}

// ListUserByUsername mocks base method
func (m *MockUserRepo) ListUserByUsername(arg0 context.Context, arg1 string, arg2 app.Page) ([]app.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserByUsername", arg0, arg1, arg2)
	ret0, _ := ret[0].([]app.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUserByUsername indicates an expected call of ListUserByUsername
func (mr *MockUserRepoMockRecorder) ListUserByUsername(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserByUsername", reflect.TypeOf((*MockUserRepo)(nil).ListUserByUsername), arg0, arg1, arg2)
}

// MockSessionRepo is a mock of SessionRepo interface
type MockSessionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepoMockRecorder
}

// MockSessionRepoMockRecorder is the mock recorder for MockSessionRepo
type MockSessionRepoMockRecorder struct {
	mock *MockSessionRepo
}

// NewMockSessionRepo creates a new mock instance
func NewMockSessionRepo(ctrl *gomock.Controller) *MockSessionRepo {
	mock := &MockSessionRepo{ctrl: ctrl}
	mock.recorder = &MockSessionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessionRepo) EXPECT() *MockSessionRepoMockRecorder {
	return m.recorder
}

// SaveSession mocks base method
func (m *MockSessionRepo) SaveSession(arg0 context.Context, arg1 app.UserID, arg2 app.TokenID, arg3 app.Origin) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSession", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSession indicates an expected call of SaveSession
func (mr *MockSessionRepoMockRecorder) SaveSession(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSession", reflect.TypeOf((*MockSessionRepo)(nil).SaveSession), arg0, arg1, arg2, arg3)
}

// SessionByTokenID mocks base method
func (m *MockSessionRepo) SessionByTokenID(arg0 context.Context, arg1 app.TokenID) (*app.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SessionByTokenID", arg0, arg1)
	ret0, _ := ret[0].(*app.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SessionByTokenID indicates an expected call of SessionByTokenID
func (mr *MockSessionRepoMockRecorder) SessionByTokenID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SessionByTokenID", reflect.TypeOf((*MockSessionRepo)(nil).SessionByTokenID), arg0, arg1)
}

// UserByTokenID mocks base method
func (m *MockSessionRepo) UserByTokenID(arg0 context.Context, arg1 app.TokenID) (*app.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByTokenID", arg0, arg1)
	ret0, _ := ret[0].(*app.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByTokenID indicates an expected call of UserByTokenID
func (mr *MockSessionRepoMockRecorder) UserByTokenID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByTokenID", reflect.TypeOf((*MockSessionRepo)(nil).UserByTokenID), arg0, arg1)
}

// DeleteSession mocks base method
func (m *MockSessionRepo) DeleteSession(arg0 context.Context, arg1 app.TokenID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession
func (mr *MockSessionRepoMockRecorder) DeleteSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockSessionRepo)(nil).DeleteSession), arg0, arg1)
}

// MockCodeRepo is a mock of CodeRepo interface
type MockCodeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCodeRepoMockRecorder
}

// MockCodeRepoMockRecorder is the mock recorder for MockCodeRepo
type MockCodeRepoMockRecorder struct {
	mock *MockCodeRepo
}

// NewMockCodeRepo creates a new mock instance
func NewMockCodeRepo(ctrl *gomock.Controller) *MockCodeRepo {
	mock := &MockCodeRepo{ctrl: ctrl}
	mock.recorder = &MockCodeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCodeRepo) EXPECT() *MockCodeRepoMockRecorder {
	return m.recorder
}

// SaveCode mocks base method
func (m *MockCodeRepo) SaveCode(ctx context.Context, id app.UserID, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCode", ctx, id, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCode indicates an expected call of SaveCode
func (mr *MockCodeRepoMockRecorder) SaveCode(ctx, id, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCode", reflect.TypeOf((*MockCodeRepo)(nil).SaveCode), ctx, id, code)
}

// UserIDByCode mocks base method
func (m *MockCodeRepo) UserIDByCode(ctx context.Context, code string) (app.UserID, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserIDByCode", ctx, code)
	ret0, _ := ret[0].(app.UserID)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UserIDByCode indicates an expected call of UserIDByCode
func (mr *MockCodeRepoMockRecorder) UserIDByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserIDByCode", reflect.TypeOf((*MockCodeRepo)(nil).UserIDByCode), ctx, code)
}

// Code mocks base method
func (m *MockCodeRepo) Code(ctx context.Context, id app.UserID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Code", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Code indicates an expected call of Code
func (mr *MockCodeRepoMockRecorder) Code(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Code", reflect.TypeOf((*MockCodeRepo)(nil).Code), ctx, id)
}

// MockCode is a mock of Code interface
type MockCode struct {
	ctrl     *gomock.Controller
	recorder *MockCodeMockRecorder
}

// MockCodeMockRecorder is the mock recorder for MockCode
type MockCodeMockRecorder struct {
	mock *MockCode
}

// NewMockCode creates a new mock instance
func NewMockCode(ctrl *gomock.Controller) *MockCode {
	mock := &MockCode{ctrl: ctrl}
	mock.recorder = &MockCodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCode) EXPECT() *MockCodeMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockCode) Generate(length int) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", length)
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate
func (mr *MockCodeMockRecorder) Generate(length interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockCode)(nil).Generate), length)
}

// MockPassword is a mock of Password interface
type MockPassword struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordMockRecorder
}

// MockPasswordMockRecorder is the mock recorder for MockPassword
type MockPasswordMockRecorder struct {
	mock *MockPassword
}

// NewMockPassword creates a new mock instance
func NewMockPassword(ctrl *gomock.Controller) *MockPassword {
	mock := &MockPassword{ctrl: ctrl}
	mock.recorder = &MockPasswordMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPassword) EXPECT() *MockPasswordMockRecorder {
	return m.recorder
}

// Hashing mocks base method
func (m *MockPassword) Hashing(password string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hashing", password)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hashing indicates an expected call of Hashing
func (mr *MockPasswordMockRecorder) Hashing(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hashing", reflect.TypeOf((*MockPassword)(nil).Hashing), password)
}

// Compare mocks base method
func (m *MockPassword) Compare(hashedPassword, password []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", hashedPassword, password)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Compare indicates an expected call of Compare
func (mr *MockPasswordMockRecorder) Compare(hashedPassword, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockPassword)(nil).Compare), hashedPassword, password)
}

// MockAuth is a mock of Auth interface
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// Token mocks base method
func (m *MockAuth) Token(expired time.Duration) (app.AuthToken, app.TokenID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Token", expired)
	ret0, _ := ret[0].(app.AuthToken)
	ret1, _ := ret[1].(app.TokenID)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Token indicates an expected call of Token
func (mr *MockAuthMockRecorder) Token(expired interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Token", reflect.TypeOf((*MockAuth)(nil).Token), expired)
}

// Parse mocks base method
func (m *MockAuth) Parse(token app.AuthToken) (app.TokenID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", token)
	ret0, _ := ret[0].(app.TokenID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse
func (mr *MockAuthMockRecorder) Parse(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockAuth)(nil).Parse), token)
}

// MockOAuth is a mock of OAuth interface
type MockOAuth struct {
	ctrl     *gomock.Controller
	recorder *MockOAuthMockRecorder
}

// MockOAuthMockRecorder is the mock recorder for MockOAuth
type MockOAuthMockRecorder struct {
	mock *MockOAuth
}

// NewMockOAuth creates a new mock instance
func NewMockOAuth(ctrl *gomock.Controller) *MockOAuth {
	mock := &MockOAuth{ctrl: ctrl}
	mock.recorder = &MockOAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOAuth) EXPECT() *MockOAuthMockRecorder {
	return m.recorder
}

// Account mocks base method
func (m *MockOAuth) Account(arg0 context.Context, arg1 string) (*app.OAuthAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Account", arg0, arg1)
	ret0, _ := ret[0].(*app.OAuthAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Account indicates an expected call of Account
func (mr *MockOAuthMockRecorder) Account(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Account", reflect.TypeOf((*MockOAuth)(nil).Account), arg0, arg1)
}
