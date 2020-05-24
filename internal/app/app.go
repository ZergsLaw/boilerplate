// Package app contains all logic of the project, all business tasks,
// manages all other modules.
package app

import (
	"errors"
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
	ErrNotValidCode              = errors.New("code not equal")
)

type (
	// App implements the business logic.
	App interface {
		UserApp
	}
	// Page for search in repo.
	Page struct {
		Limit  int // > 0
		Offset int // >= 0
	}
	// Application implements interface App.
	Application struct {
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

// Config for build project.
type Config struct {
	UserRepo     UserRepo
	SessionRepo  SessionRepo
	CodeRepo     CodeRepo
	Password     Password
	Auth         Auth
	Wal          WAL
	Notification Notification
	Code         Code
}

// New creates and returns new App.
func New(cfg Config) *Application {
	return &Application{
		userRepo:     cfg.UserRepo,
		sessionRepo:  cfg.SessionRepo,
		codeRepo:     cfg.CodeRepo,
		password:     cfg.Password,
		auth:         cfg.Auth,
		wal:          cfg.Wal,
		code:         cfg.Code,
		notification: cfg.Notification,
	}
}
