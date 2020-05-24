package repo

import (
	"net"
	"time"

	"github.com/jackc/pgtype"
	"github.com/zergslaw/boilerplate/internal/app"
)

type (
	userDBFormat struct {
		ID        app.UserID   `db:"id"`
		Email     string       `db:"email"`
		Username  string       `db:"username"`
		PassHash  pgtype.Bytea `db:"pass_hash"`
		CreatedAt time.Time    `db:"created_at"`
		UpdatedAt time.Time    `db:"updated_at"`
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

	codeInfoDBFormat struct {
		ID        int       `db:"id"`
		Code      string    `db:"code"`
		Email     string    `db:"email"`
		CreatedAt time.Time `db:"created_at"`
	}

	taskNotificationDBFormat struct {
		ID    int    `db:"id"`
		Email string `db:"email"`
		Kind  string `db:"kind"`
	}
)

func (val *userDBFormat) toAppFormat() *app.User {
	return &app.User{
		ID:        val.ID,
		Email:     val.Email,
		Name:      val.Username,
		PassHash:  val.PassHash.Bytes,
		CreatedAt: val.CreatedAt,
		UpdatedAt: val.UpdatedAt,
	}
}

func (val *sessionDBFormat) toAppFormat() *app.Session {
	return &app.Session{
		Origin: app.Origin{
			IP:        val.IP.IPNet.IP,
			UserAgent: val.UserAgent,
		},
		ID:      val.ID,
		TokenID: app.TokenID(val.TokenID),
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

func (val *taskNotificationDBFormat) toAppFormat() *app.TaskNotification {
	kind := app.Welcome
	switch val.Kind {
	case app.ChangeEmail.String():
		kind = app.ChangeEmail
	case app.PassRecovery.String():
		kind = app.PassRecovery
	}

	return &app.TaskNotification{
		ID:    val.ID,
		Email: val.Email,
		Kind:  kind,
	}
}

func (val *codeInfoDBFormat) toAppFormat() *app.CodeInfo {
	return &app.CodeInfo{
		Code:      val.Code,
		Email:     val.Email,
		CreatedAt: val.CreatedAt,
	}
}
