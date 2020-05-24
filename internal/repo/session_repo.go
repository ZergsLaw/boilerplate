package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"

	"github.com/jmoiron/sqlx"

	"github.com/zergslaw/boilerplate/internal/app"
)

// SaveSession need for implements app.SessionRepo.
func (repo *Repo) SaveSession(ctx context.Context, userID app.UserID, tokenID app.TokenID, origin app.Origin) error {
	return repo.db.Do(func(db *sqlx.DB) error {
		const query = `INSERT INTO sessions (user_id, token_id, ip, user_agent) VALUES (:user_id,:token_id,:ip,:user_agent)`
		type args struct {
			UserID    app.UserID   `db:"user_id"`
			TokenID   app.TokenID  `db:"token_id"`
			IP        *pgtype.Inet `db:"ip"`
			UserAgent string       `db:"user_agent"`
		}

		inet, err := inet(origin.IP)
		if err != nil {
			return fmt.Errorf("inet: %w", err)
		}

		_, err = db.NamedExecContext(ctx, query, args{
			UserID:    userID,
			TokenID:   tokenID,
			IP:        inet,
			UserAgent: origin.UserAgent,
		})
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		return nil
	})
}

// SessionByTokenID need for implements app.SessionRepo.
func (repo *Repo) SessionByTokenID(ctx context.Context, tokenID app.TokenID) (session *app.Session, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM sessions WHERE token_id = $1 AND is_logout = false`

		s := &sessionDBFormat{}
		err = db.GetContext(ctx, s, query, tokenID)
		if err != nil {
			return err
		}

		session = s.toAppFormat()
		return nil
	})
	return
}

// UserByTokenID need for implements app.UserRepo.
func (repo *Repo) UserByTokenID(ctx context.Context, tokenID app.TokenID) (user *app.User, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT users.id, users.email, users.username, users.pass_hash, users.created_at, users.updated_at
		FROM users LEFT JOIN sessions ON sessions.user_id = users.id WHERE sessions.token_id = $1
		AND sessions.is_logout = false`

		u := &userDBFormat{}
		err = db.GetContext(ctx, u, query, tokenID)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// DeleteSession need for implements app.SessionRepo.
func (repo *Repo) DeleteSession(ctx context.Context, tokenID app.TokenID) error {
	return repo.db.Do(func(db *sqlx.DB) error {
		const query = `UPDATE sessions SET is_logout = true WHERE token_id = :token_id`
		type args struct {
			TokenID app.TokenID `db:"token_id"`
		}

		_, err := db.NamedExecContext(ctx, query, args{TokenID: tokenID})
		return err
	})
}
