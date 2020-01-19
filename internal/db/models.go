package db

import (
	"net"
	"time"

	"github.com/jackc/pgtype"
	"github.com/zergslaw/users/internal/app"
)

type (
	userDBFormat struct {
		ID        app.UserID   `db:"id"`
		Email     string       `db:"email"`
		Username  string       `db:"username"`
		PassHash  pgtype.Bytea `db:"pass_hash"`
		CreatedAt time.Time    `db:"created_at"`
		UpdatedAt time.Time    `db:"updated_at"`

		// Need to get the number of tracks conveniently found.
		Total int `db:"total"`
	}

	sessionDBFormat struct {
		ID        app.SessionID `db:"id"`
		UserID    app.UserID    `db:"user_id"`
		TokenID   app.AuthToken `db:"token_id"`
		IP        *pgtype.Inet  `db:"ip"`
		UserAgent string        `db:"user_agent"`
		IsLogout  bool          `db:"is_logout"`
		CreatedAt time.Time     `db:"created_at"`
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
