package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/zergslaw/users/internal/app"
)

// CreateUser need for implements app.Repo.
func (repo *Repo) CreateUser(ctx context.Context, newUser app.User) (userID app.UserID, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `INSERT INTO users (username, email, pass_hash) VALUES (:username, :email, :pass_hash)`
		type arg struct {
			Username string       `db:"username"`
			Email    string       `db:"email"`
			PassHash pgtype.Bytea `db:"pass_hash"`
		}

		_, err := db.NamedExecContext(ctx, query, arg{
			Username: newUser.Username,
			Email:    newUser.Email,
			PassHash: pgtype.Bytea{
				Bytes:  newUser.PassHash,
				Status: pgtype.Present,
			},
		})
		switch {
		case pqErrConflictIn(err, constraintEmail):
			return app.ErrEmailExist
		case pqErrConflictIn(err, constraintUsername):
			return app.ErrUsernameExist
		case err != nil:
			return fmt.Errorf("create user: %w", err)
		}

		user, err := repo.UserByEmail(ctx, newUser.Email)
		if err != nil {
			return err
		}

		userID = user.ID
		return nil
	})
	return // nolint:nakedret
}

// DeleteUser need for implements app.Repo.
func (repo *Repo) DeleteUser(ctx context.Context, userID app.UserID) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `DELETE FROM users WHERE id = $1`
		_, err := db.ExecContext(ctx, query, userID)

		return err
	})
}

// UpdateUsername need for implements app.Repo.
func (repo *Repo) UpdateUsername(ctx context.Context, userID app.UserID, username string) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `UPDATE users SET username = :username, updated_at = now() WHERE id = :id`
		type arg struct {
			ID       app.UserID `db:"id"`
			Username string     `db:"username"`
		}

		_, err := db.NamedExecContext(ctx, query, arg{
			ID:       userID,
			Username: username,
		})
		switch {
		case pqErrConflictIn(err, constraintUsername):
			return app.ErrUsernameExist
		case err != nil:
			return err
		}

		return nil
	})
}

// UpdateEmail need for implements app.Repo.
func (repo *Repo) UpdateEmail(ctx context.Context, userID app.UserID, email string) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `UPDATE users SET email = :email, updated_at = now() WHERE id = :id`
		type arg struct {
			ID    app.UserID `db:"id"`
			Email string     `db:"email"`
		}

		_, err := db.NamedExecContext(ctx, query, arg{
			ID:    userID,
			Email: email,
		})
		switch {
		case pqErrConflictIn(err, constraintEmail):
			return app.ErrEmailExist
		case err != nil:
			return err
		}

		return nil
	})
}

// UpdatePassword need for implements app.Repo.
func (repo *Repo) UpdatePassword(ctx context.Context, userID app.UserID, passHash []byte) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `UPDATE users SET pass_hash = :pass_hash, updated_at = now() WHERE id = :id`
		type arg struct {
			ID       app.UserID   `db:"id"`
			PassHash pgtype.Bytea `db:"pass_hash"`
		}

		_, err := db.NamedExecContext(ctx, query, arg{
			ID: userID,
			PassHash: pgtype.Bytea{
				Bytes:  passHash,
				Status: pgtype.Present,
			},
		})

		return err
	})
}

// UserByID need for implements app.Repo.
func (repo *Repo) UserByID(ctx context.Context, userID app.UserID) (user *app.User, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE id = $1`

		u := &userDBFormat{}
		err := get(ctx, db, u, query, userID)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// UserByTokenID need for implements app.Repo.
func (repo *Repo) UserByTokenID(ctx context.Context, token app.TokenID) (user *app.User, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT users.id, users.email, users.username, users.pass_hash, users.created_at, users.updated_at
		FROM users LEFT JOIN sessions ON sessions.user_id = users.id WHERE sessions.token_id = $1
		AND sessions.is_logout = false`

		u := &userDBFormat{}
		err := get(ctx, db, u, query, token)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// UserByEmail need for implements app.Repo.
func (repo *Repo) UserByEmail(ctx context.Context, email string) (user *app.User, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE email = $1`

		u := &userDBFormat{}
		err := get(ctx, db, u, query, email)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// UserByUsername need for implements app.Repo.
func (repo *Repo) UserByUsername(ctx context.Context, username string) (user *app.User, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE username = $1`

		u := &userDBFormat{}
		err := get(ctx, db, u, query, username)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// ListUserByUsername need for implements app.Repo.
func (repo *Repo) ListUserByUsername(ctx context.Context, username string, page app.Page) (users []app.User, total int, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT *, count(*) OVER() AS total FROM users WHERE username LIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3 `

		var items []userDBFormat
		err := db.SelectContext(ctx, &items, query, "%"+username+"%", page.Limit, page.Offset)
		if err != nil {
			return fmt.Errorf("db select: %w", err)
		}

		users = make([]app.User, len(items))
		for i := range items {
			users[i] = *items[i].toAppFormat()
			total = items[i].Total
		}

		return nil
	})
	return
}

// SaveSession need for implements app.Repo.
func (repo *Repo) SaveSession(ctx context.Context, userID app.UserID, tokenID app.TokenID, origin app.Origin) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `INSERT INTO sessions (user_id, token_id, ip, user_agent) VALUES (:user_id, :token_id, :ip, :user_agent)`

		type arg struct {
			UserID    app.UserID   `db:"user_id"`
			TokenID   app.TokenID  `db:"token_id"`
			IP        *pgtype.Inet `db:"ip"`
			UserAgent string       `db:"user_agent"`
		}

		inet, err := inet(origin.IP)
		if err != nil {
			return fmt.Errorf("inet: %w", err)
		}

		_, err = db.NamedExecContext(ctx, query, arg{
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

// SessionByTokenID need for implements app.Repo.
func (repo *Repo) SessionByTokenID(ctx context.Context, tokenID app.TokenID) (session *app.Session, err error) {
	err = repo.execFunc(func(db *sqlx.DB) error {
		const query = `SELECT * FROM sessions WHERE token_id = $1 AND is_logout = false`

		sessionFromDB := &sessionDBFormat{}
		err := get(ctx, db, sessionFromDB, query, tokenID)
		if err != nil {
			return err
		}

		session = sessionFromDB.toAppFormat()
		return nil
	})
	return
}

// DeleteSession need for implements app.Repo.
func (repo *Repo) DeleteSession(ctx context.Context, tokenID app.TokenID) error {
	return repo.execFunc(func(db *sqlx.DB) error {
		const query = `UPDATE sessions SET is_logout = true WHERE token_id = $1`
		_, err := db.ExecContext(ctx, query, tokenID)

		return err
	})
}

func get(ctx context.Context, db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := db.GetContext(ctx, dest, query, args...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return app.ErrNotFound
	case err != nil:
		return fmt.Errorf("db get: %w", err)
	}

	return nil
}
