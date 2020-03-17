// Package app contains all logic of the project, all business tasks,
// manages all other modules.
package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// Errors.
var (
	ErrEmailExist                = errors.New("email exist")
	ErrUsernameExist             = errors.New("username exist")
	ErrNotFound                  = errors.New("not found")
	ErrNotValidPassword          = errors.New("not valid password")
	ErrInvalidToken              = errors.New("not valid auth")
	ErrExpiredToken              = errors.New("auth is expired")
	ErrUsernameNeedDifferentiate = errors.New("username need to differentiate")
	ErrEmailNeedDifferentiate    = errors.New("email need to differentiate")
	ErrNotUnknownKindTask        = errors.New("unknown task kind")
	ErrCodeExpired               = errors.New("code is expired")
)

type (
	// UserRepo interface for user data repository.
	UserRepo interface {
		// CreateUser adds to the new user in repository.
		// This method is also required to create a notifying hoard.
		// Errors: ErrEmailExist, ErrUsernameExist, unknown.
		CreateUser(context.Context, User) (UserID, error)
		// DeleteUser removes user from repository.
		// Errors: unknown.
		DeleteUser(context.Context, UserID) error
		// UpdateUsername changes username if he's not busy.
		// Errors: ErrUsernameExist, unknown.
		UpdateUsername(context.Context, UserID, string) error
		// UpdateEmail changes email if he's not busy.
		// Errors: ErrEmailExist, unknown.
		UpdateEmail(context.Context, UserID, string) error
		// UpdatePassword changes password.
		// Resets all codes to reset the password.
		// Errors: unknown.
		UpdatePassword(context.Context, UserID, []byte) error
		// UserByID returning user info by id.
		// Errors: ErrNotFound, unknown.
		UserByID(context.Context, UserID) (*User, error)
		// UserByAuthToken returning user info by authToken.
		// Errors: ErrNotFound, unknown.
		UserByTokenID(context.Context, TokenID) (*User, error)
		// UserByEmail returning user info by email.
		// This method is also required to create a notifying hoard.
		// Errors: ErrNotFound, unknown.
		UserByEmail(context.Context, string) (*User, error)
		// UserByUsername returning user info by id.
		// Errors: ErrNotFound, unknown.
		UserByUsername(context.Context, string) (*User, error)
		// ListUserByUsername returning list user info.
		// Errors: unknown.
		ListUserByUsername(context.Context, string, Page) ([]User, int, error)
	}
	// SessionRepo interface for session data repository.
	SessionRepo interface {
		// SaveSession saves the new user Session in a database.
		// Errors: unknown.
		SaveSession(context.Context, UserID, TokenID, Origin) error
		// Session returns user Session.
		// Errors: ErrNotFound, unknown.
		SessionByTokenID(context.Context, TokenID) (*Session, error)
		// DeleteSession removes user Session.
		// Errors: unknown.
		DeleteSession(context.Context, TokenID) error
	}
	// CodeRepo interface for recover code repository.
	CodeRepo interface {
		// SaveCode the code to restore the password to the repository.
		// Removes all recovery codes from this email before adding a new one.
		// Creates a task to send the recovery code to the user's mail.
		// Errors: unknown.
		SaveCode(ctx context.Context, id UserID, code string) error
		// UserID returns user id by recovery code.
		// Errors: ErrNotFound, unknown.
		UserID(ctx context.Context, code string) (userID UserID, createAt time.Time, err error)
		// Code returns recovery code for recovery password by user id.
		// Errors: ErrNotFound, unknown.
		Code(ctx context.Context, id UserID) (code string, err error)
	}
	// WAL module returning tasks and also closing them.
	WAL interface {
		// NotificationTask returns the earliest task that has not been completed.
		// Errors: ErrNotFound, unknown.
		NotificationTask(ctx context.Context) (task *TaskNotification, err error)
		// DeleteTaskNotification removes the task performed.
		// Errors: unknown.
		DeleteTaskNotification(ctx context.Context, id int) error
	}
	// Notification module for working with alerts for registered users.
	Notification interface {
		// NotificationTask must accept the parameter contact to whom the notification will be sent.
		// At the moment, the guarantee of message delivery lies on this module, it is possible to
		// transfer it to the app.
		Notification(contact string, msg Message) error
	}
	// Code module for generate random code.
	Code interface {
		// Generate random code of a specified length.
		Generate(length int) string
	}
	// Password module responsible for working with passwords.
	Password interface {
		// Hashing returns the hashed version of the password.
		// Errors: unknown.
		Hashing(password string) ([]byte, error)
		// Compare compares two passwords for matches.
		Compare(hashedPassword []byte, password []byte) bool
	}
	// Auth module is responsible for working with authorization tokens.
	Auth interface {
		// Token generates an authorization auth with a specified lifetime,
		// and can also use the UserID if necessary.
		// Errors: unknown.
		Token(expired time.Duration) (AuthToken, TokenID, error)
		// Parse and validates the auth and checks that it's expired.
		// Errors: ErrInvalidToken, ErrExpiredToken, unknown.
		Parse(token AuthToken) (TokenID, error)
	}
	// OAuth module responsible for working with social network.
	// TODO Implements.
	OAuth interface {
		// Account converts an authorization code into user information.
		// Errors: unknown.
		Account(context.Context, string) (*OAuthAccount, error)
	}
	// App implements the business logic.
	App interface {
		// VerificationEmail check if the user is registered with this email.
		// Errors: ErrEmailExist, unknown.
		VerificationEmail(ctx context.Context, email string) error
		// VerificationUsername check if the user is registered with this username.
		// Errors: ErrUsernameExist, unknown.
		VerificationUsername(ctx context.Context, username string) error
		// Login authorizes the user to the system.
		// Errors: ErrNotFound, ErrNotValidPassword, unknown.
		Login(ctx context.Context, email, password string, origin Origin) (*User, AuthToken, error)
		// Logout remove user Session.
		// Errors: unknown.
		Logout(context.Context, AuthUser) error
		// CreateUser creates a new user to the system, the password is hashed with bcrypt.
		// Errors: ErrEmailExist, ErrUsernameExist, unknown.
		CreateUser(ctx context.Context, email, username, password string, origin Origin) (*User, AuthToken, error)
		// DeleteUser deleting user profile.
		// Errors: unknown.
		DeleteUser(context.Context, AuthUser) error
		// User returning user profile.
		// Errors: ErrNotFound, unknown.
		User(context.Context, AuthUser, UserID) (*User, error)
		// UserByAuthToken returns user by authToken.
		// Errors: ErrNotFound, unknown.
		UserByAuthToken(ctx context.Context, token AuthToken) (*AuthUser, error)
		// UpdateUsername refresh the username.
		// Errors: ErrUsernameExist, ErrUsernameNeedDifferentiate, unknown.
		UpdateUsername(context.Context, AuthUser, string) error
		// UpdateEmail refresh the email.
		// Errors: ErrEmailExist, ErrEmailNeedDifferentiate, unknown.
		UpdateEmail(context.Context, AuthUser, string) error
		// UpdateUsername refresh user password.
		// Errors: ErrNotValidPassword, unknown.
		UpdatePassword(ctx context.Context, authUser AuthUser, oldPass, newPass string) error
		// ListUserByUsername returns list user by username.
		// Errors: unknown.
		ListUserByUsername(context.Context, AuthUser, string, Page) ([]User, int, error)
		// CreateRecoveryCode creates and sends a password recovery code to the user's email.
		// Errors: ErrNotFound, unknown.
		CreateRecoveryCode(ctx context.Context, email string) error
		// RecoveryPassword replaces the password with a new one from the user who owns this recovery code.
		// Errors: ErrCodeExpired, ErrNotFound, unknown.
		RecoveryPassword(ctx context.Context, code, newPassword string) error
		// StartWALNotification starts the task of notifying users.
		StartWALNotification(ctx context.Context) error
	}
	// UserID contains user id.
	UserID int
	// SessionID contains Session id.
	SessionID int
	// SocialID contains id from social network with OAuth.
	SocialID string
	// AuthToken authorization auth.
	AuthToken string
	// TokenID contains auth id.
	TokenID string
	// TaskNotification contains information to perform the task of notifying the user.
	TaskNotification struct {
		ID     int
		UserID UserID
		Kind   MessageKind
	}
	// MessageKind selects the type of message to be sent.
	MessageKind int
	// Message contains sent info.
	Message struct {
		Kind    MessageKind
		Content string
	}
	// Page for search users in repo.
	Page struct {
		Limit  int // > 0
		Offset int // >= 0
	}
	// Origin information about req user.
	Origin struct {
		IP        net.IP
		UserAgent string
	}
	// OAuthAccount user information from the social network.
	OAuthAccount struct {
		ID       string
		Email    string
		Username string
	}
	// Session contains user Session information.
	Session struct {
		Origin
		ID      SessionID
		TokenID TokenID
	}
	// User contains user information.
	User struct {
		ID       UserID
		Email    string
		Username string
		PassHash []byte

		CreatedAt time.Time
		UpdatedAt time.Time
	}
	// AuthUser contains auth information.
	AuthUser struct {
		User
		Session Session
	}
	// app implements interface App.
	app struct {
		userRepo     UserRepo
		sessionRepo  SessionRepo
		codeRepo     CodeRepo
		password     Password
		auth         Auth
		wal          WAL
		notification Notification
		code         Code
	}
)

// Message enums.
const (
	Welcome MessageKind = iota + 1
	ChangeEmail
	PassRecovery
)

func (m MessageKind) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case ChangeEmail:
		return "change email"
	case PassRecovery:
		return "password recovery"
	default:
		panic(fmt.Sprintf("unknown kind: %d", m))
	}
}

// New creates and returns new App.
func New(userRepo UserRepo, sessionRepo SessionRepo, codeRepo CodeRepo,
	password Password, auth Auth, wal WAL, notification Notification,
	code Code) App {
	return &app{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		codeRepo:     codeRepo,
		password:     password,
		auth:         auth,
		wal:          wal,
		code:         code,
		notification: notification,
	}
}
