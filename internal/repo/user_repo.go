package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/zergslaw/boilerplate/internal/app"
)

// CreateUser need for implements app.UserRepo.
func (repo *Repo) CreateUser(ctx context.Context, newUser app.User, task app.TaskNotification) (userID app.UserID, err error) {
	err = repo.db.Tx(ctx, func(tx *sqlx.Tx) error {
		const query = `INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id`

		hash := pgtype.Bytea{
			Bytes:  newUser.PassHash,
			Status: pgtype.Present,
		}

		err = tx.QueryRowxContext(ctx, query, newUser.Name, newUser.Email, hash).Scan(&userID)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return createTaskNotification(ctx, tx, task)
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// DeleteUser need for implements app.UserRepo.
func (repo *Repo) DeleteUser(ctx context.Context, userID app.UserID) error {
	return repo.db.Do(func(db *sqlx.DB) error {
		const query = `DELETE FROM users WHERE id = :id`
		type args struct {
			ID app.UserID `db:"id"`
		}

		_, err := db.NamedExecContext(ctx, query, args{
			ID: userID,
		})

		return err
	})
}

// UpdateUsername need for implements app.UserRepo.
func (repo *Repo) UpdateUsername(ctx context.Context, userID app.UserID, username string) error {
	return repo.db.Do(func(db *sqlx.DB) error {
		const query = `UPDATE users SET username = :username, updated_at = now() WHERE id = :id`
		type args struct {
			Username string     `db:"username"`
			ID       app.UserID `db:"id"`
		}

		_, err := db.NamedExecContext(ctx, query, args{
			Username: username,
			ID:       userID,
		})

		return err
	})
}

// UpdateEmail need for implements app.UserRepo.
func (repo *Repo) UpdateEmail(ctx context.Context, userID app.UserID, email string, task app.TaskNotification) error {
	return repo.db.Tx(ctx, func(tx *sqlx.Tx) error {
		const query = `UPDATE users SET email = :email, updated_at = now() WHERE id = :id`
		type args struct {
			Email string     `db:"email"`
			ID    app.UserID `db:"id"`
		}

		_, err := tx.NamedExecContext(ctx, query, args{
			Email: email,
			ID:    userID,
		})
		if err != nil {
			return fmt.Errorf("update email: %w", err)
		}

		return createTaskNotification(ctx, tx, task)
	})
}

// UpdatePassword need for implements app.UserRepo.
func (repo *Repo) UpdatePassword(ctx context.Context, userID app.UserID, passHash []byte) error {
	return repo.db.Tx(ctx, func(tx *sqlx.Tx) error {
		const query = `UPDATE users SET pass_hash = $1, updated_at = now() WHERE id = $2 RETURNING email`
		hash := pgtype.Bytea{
			Bytes:  passHash,
			Status: pgtype.Present,
		}

		userEmail := ""
		err := tx.QueryRowContext(ctx, query, hash, userID).Scan(&userEmail)
		if err != nil {
			return fmt.Errorf("update pass: %w", err)
		}

		return cleanRecoveryCodes(ctx, tx, userEmail)
	})
}

// UserByID need for implements app.UserRepo.
func (repo *Repo) UserByID(ctx context.Context, userID app.UserID) (user *app.User, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE id = $1`

		u := &userDBFormat{}
		err = db.GetContext(ctx, u, query, userID)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// UserByEmail need for implements app.UserRepo.
func (repo *Repo) UserByEmail(ctx context.Context, email string) (user *app.User, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE email = $1`

		u := &userDBFormat{}
		err = db.GetContext(ctx, u, query, email)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// UserByUsername need for implements app.UserRepo.
func (repo *Repo) UserByUsername(ctx context.Context, username string) (user *app.User, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE username = $1`

		u := &userDBFormat{}
		err = db.GetContext(ctx, u, query, username)
		if err != nil {
			return err
		}

		user = u.toAppFormat()
		return nil
	})
	return
}

// ListUserByUsername need for implements app.UserRepo.
func (repo *Repo) ListUserByUsername(ctx context.Context, username string, page app.Page) (users []app.User, total int, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM users WHERE username LIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

		res := make([]userDBFormat, 0, page.Limit)
		err = db.SelectContext(ctx, &res, query, "%"+username+"%", page.Limit, page.Offset)
		if err != nil {
			return fmt.Errorf("select: %w", err)
		}

		const getTotal = `SELECT count(*) OVER() AS total FROM users WHERE username LIKE $1`
		err = db.GetContext(ctx, &total, getTotal, "%"+username+"%")
		if err != nil {
			return fmt.Errorf("get total: %w", err)
		}

		users = make([]app.User, len(res))
		for i := range res {
			users[i] = *res[i].toAppFormat()
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
