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
)

type (
	// UserRepo interface for data repository.
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
		// Errors: unknown.
		UpdatePassword(context.Context, UserID, []byte) error
		// UserByID returning user info by id.
		// Errors: ErrNotFound, unknown.
		UserByID(context.Context, UserID) (*User, error)
		// UserByAuthToken returning user info by authToken.
		// Errors: ErrNotFound, unknown.
		UserByTokenID(context.Context, TokenID) (*User, error)
		// UserByEmail returning user info by email.
		// Errors: ErrNotFound, unknown.
		UserByEmail(context.Context, string) (*User, error)
		// UserByUsername returning user info by id.
		// Errors: ErrNotFound, unknown.
		UserByUsername(context.Context, string) (*User, error)
		// ListUserByUsername returning list user info.
		// Errors: unknown.
		ListUserByUsername(context.Context, string, Page) ([]User, int, error)
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
	// WAL module returning tasks and also closing them.
	WAL interface {
		// NotificationTask returns the earliest task that has not been completed.
		// Errors: ErrNoTasks, unknown.
		NotificationTask(ctx context.Context) (task *TaskNotification, err error)
		// DeleteTaskNotification removes the task performed.
		// Errors: ErrNotFound, unknown.
		DeleteTaskNotification(ctx context.Context, id int) error
	}
	// Notification module for working with alerts for registered users.
	Notification interface {
		// NotificationTask must accept the parameter contact to whom the notification will be sent.
		// At the moment, the guarantee of message delivery lies on this module, it is possible to
		// transfer it to the app.
		Notification(contact string, msg Message) error
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
		ID    int
		Email string
		Kind  Message
	}
	// Message selects the type of message to be sent..
	Message int
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
		repo         UserRepo
		password     Password
		auth         Auth
		wal          WAL
		notification Notification
	}
)

// Message enums.
const (
	Welcome Message = iota + 1
	ChangeEmail
)

func (m Message) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case ChangeEmail:
		return "change email"
	default:
		panic(fmt.Sprintf("unknown kind: %d", m))
	}
}

// New creates and returns new App.
func New(repo UserRepo, password Password, auth Auth, wal WAL, notification Notification) App {
	return &app{
		repo:         repo,
		password:     password,
		auth:         auth,
		wal:          wal,
		notification: notification,
	}
}
