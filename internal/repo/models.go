package repo

import (
	"net"
	"time"

	"github.com/jackc/pgtype"
	"github.com/zergslaw/boilerplate/internal/app"
)

type (
	userDBFormat struct {
		ID        app.UserID
		Email     string
		Username  string
		PassHash  pgtype.Bytea
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	sessionDBFormat struct {
		ID        app.SessionID
		UserID    app.UserID
		TokenID   app.AuthToken
		IP        *pgtype.Inet
		UserAgent string
		IsLogout  bool
		CreatedAt time.Time
	}

	taskNotificationDBFormat struct {
		ID     int
		UserID app.UserID
		Kind   string
	}
)

func (u *userDBFormat) toAppFormat() *app.User {
	return &app.User{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		PassHash:  u.PassHash.Bytes,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (s *sessionDBFormat) toAppFormat() *app.Session {
	return &app.Session{
		Origin: app.Origin{
			IP:        s.IP.IPNet.IP,
			UserAgent: s.UserAgent,
		},
		ID:      s.ID,
		TokenID: app.TokenID(s.TokenID),
	}
}

func inet(ip net.IP) (*pgtype.Inet, error) {
	inet := &pgtype.Inet{}
	if ip == nil || ip.IsUnspecified() {
		err := inet.Set(nil)
		if err != nil {
			return nil, err
		}
	} else {
		err := inet.Set(ip)
		if err != nil {
			return nil, err
		}
	}

	return inet, nil
}

func (s *taskNotificationDBFormat) toAppFormat() *app.TaskNotification {
	kind := app.Welcome
	switch s.Kind {
	case app.ChangeEmail.String():
		kind = app.ChangeEmail
	case app.PassRecovery.String():
		kind = app.PassRecovery
	}

	return &app.TaskNotification{
		ID:     s.ID,
		UserID: s.UserID,
		Kind:   kind,
	}
}
