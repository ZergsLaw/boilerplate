package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/zergslaw/boilerplate/internal/app"
)

// CreateUser need for implements app.UserRepo.
func (repo *Repo) CreateUser(ctx context.Context, newUser app.User) (userID app.UserID, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id`

		hash := pgtype.Bytea{
			Bytes:  newUser.PassHash,
			Status: pgtype.Present,
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		err = tx.QueryRowContext(ctx, query, newUser.Username, newUser.Email, hash).Scan(&userID)
		switch {
		case pqErrConflictIn(err, constraintEmail):
			return rollback(tx, app.ErrEmailExist)
		case pqErrConflictIn(err, constraintUsername):
			return rollback(tx, app.ErrUsernameExist)
		case err != nil:
			return rollback(tx, fmt.Errorf("create user: %w", err))
		}

		err = createTaskNotification(ctx, tx, userID, app.Welcome)
		if err != nil {
			return rollback(tx, err)
		}

		return tx.Commit()
	})
	return userID, err
}

func rollback(tx *sql.Tx, err error) error {
	errRollback := tx.Rollback()
	if errRollback != nil {
		err = fmt.Errorf("%w, err rollback: %s", err, errRollback)
	}

	return err
}

// DeleteUser need for implements app.UserRepo.
func (repo *Repo) DeleteUser(ctx context.Context, userID app.UserID) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `DELETE FROM users WHERE id = $1`
		_, err := db.ExecContext(ctx, query, userID)

		return err
	})
}

// UpdateUsername need for implements app.UserRepo.
func (repo *Repo) UpdateUsername(ctx context.Context, userID app.UserID, username string) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `UPDATE users SET username = $1, updated_at = now() WHERE id = $2`

		_, err := db.ExecContext(ctx, query, username, userID)
		switch {
		case pqErrConflictIn(err, constraintUsername):
			return app.ErrUsernameExist
		case err != nil:
			return err
		}

		return nil
	})
}

// UpdateEmail need for implements app.UserRepo.
func (repo *Repo) UpdateEmail(ctx context.Context, userID app.UserID, email string) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `UPDATE users SET email = $1, updated_at = now() WHERE id = $2`

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		_, err = tx.ExecContext(ctx, query, email, userID)
		switch {
		case pqErrConflictIn(err, constraintEmail):
			return rollback(tx, app.ErrEmailExist)
		case err != nil:
			return rollback(tx, err)
		}

		err = createTaskNotification(ctx, tx, userID, app.ChangeEmail)
		if err != nil {
			return rollback(tx, err)
		}

		return tx.Commit()
	})
}

// UpdatePassword need for implements app.UserRepo.
func (repo *Repo) UpdatePassword(ctx context.Context, userID app.UserID, passHash []byte) error {
	return repo.execFunc(func(db *sql.DB) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		const query = `UPDATE users SET pass_hash = $1, updated_at = now() WHERE id = $2 RETURNING email`

		hash := pgtype.Bytea{
			Bytes:  passHash,
			Status: pgtype.Present,
		}

		userEmail := ""
		err = tx.QueryRowContext(ctx, query, hash, userID).Scan(&userEmail)
		if err != nil {
			return rollback(tx, err)
		}

		err = cleanRecoveryCodes(ctx, tx, userID)
		if err != nil {
			return rollback(tx, err)
		}

		return tx.Commit()
	})
}

// UserByID need for implements app.UserRepo.
func (repo *Repo) UserByID(ctx context.Context, userID app.UserID) (user *app.User, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT * FROM users WHERE id = $1`

		u := &userDBFormat{}
		err = db.QueryRowContext(ctx, query, userID).Scan(
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

// UserByEmail need for implements app.UserRepo.
func (repo *Repo) UserByEmail(ctx context.Context, email string) (user *app.User, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT * FROM users WHERE email = $1`

		u := &userDBFormat{}
		err = db.QueryRowContext(ctx, query, email).Scan(
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

// UserByUsername need for implements app.UserRepo.
func (repo *Repo) UserByUsername(ctx context.Context, username string) (user *app.User, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT * FROM users WHERE username = $1`

		u := &userDBFormat{}
		err = db.QueryRowContext(ctx, query, username).Scan(
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

// ListUserByUsername need for implements app.UserRepo.
func (repo *Repo) ListUserByUsername(ctx context.Context, username string, page app.Page) (users []app.User, total int, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT *, count(*) OVER() AS total FROM users WHERE username LIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

		rows, err := db.QueryContext(ctx, query, "%"+username+"%", page.Limit, page.Offset)
		if err != nil {
			return fmt.Errorf("query context: %w", err)
		}

		if err = rows.Err(); err != nil {
			return fmt.Errorf("rows error: %w", err)
		}

		res := make([]userDBFormat, 0, page.Limit)
		for rows.Next() {
			u := userDBFormat{}
			err = rows.Scan(
				&u.ID,
				&u.Email,
				&u.Username,
				&u.PassHash,
				&u.CreatedAt,
				&u.UpdatedAt,
				&total,
			)
			if err != nil {
				return fmt.Errorf("scan user: %w", err)
			}

			res = append(res, u)
		}

		users = make([]app.User, len(res))
		for i := range res {
			users[i] = *res[i].toAppFormat()
		}

		return rows.Close()
	})
	return users, total, err
}
