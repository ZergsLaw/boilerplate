package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zergslaw/boilerplate/internal/app"
)

// SaveSession need for implements app.SessionRepo.
func (repo *Repo) SaveSession(ctx context.Context, userID app.UserID, tokenID app.TokenID, origin app.Origin) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `INSERT INTO sessions (user_id, token_id, ip, user_agent) VALUES ($1, $2, $3, $4)`

		inet, err := inet(origin.IP)
		if err != nil {
			return fmt.Errorf("inet: %w", err)
		}

		_, err = db.ExecContext(ctx, query, userID, tokenID, inet, origin.UserAgent)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		return nil
	})
}

// SessionByTokenID need for implements app.SessionRepo.
func (repo *Repo) SessionByTokenID(ctx context.Context, tokenID app.TokenID) (session *app.Session, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT * FROM sessions WHERE token_id = $1 AND is_logout = false`

		item := &sessionDBFormat{}
		err := db.QueryRowContext(ctx, query, tokenID).Scan(
			&item.ID,
			&item.UserID,
			&item.TokenID,
			&item.IP,
			&item.UserAgent,
			&item.CreatedAt,
			&item.IsLogout,
		)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return fmt.Errorf("query row: %w", err)
		}

		session = item.toAppFormat()
		return nil
	})
	return
}

// UserByTokenID need for implements app.UserRepo.
func (repo *Repo) UserByTokenID(ctx context.Context, token app.TokenID) (user *app.User, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT users.id, users.email, users.username, users.pass_hash, users.created_at, users.updated_at
		FROM users LEFT JOIN sessions ON sessions.user_id = users.id WHERE sessions.token_id = $1
		AND sessions.is_logout = false`

		u := &userDBFormat{}
		err = db.QueryRowContext(ctx, query, token).Scan(
			&u.ID,
			&u.Email,
			&u.Username,
			&u.PassHash,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// DeleteSession need for implements app.SessionRepo.
func (repo *Repo) DeleteSession(ctx context.Context, tokenID app.TokenID) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `UPDATE sessions SET is_logout = true WHERE token_id = $1`
		_, err := db.ExecContext(ctx, query, tokenID)

		return err
	})
}
