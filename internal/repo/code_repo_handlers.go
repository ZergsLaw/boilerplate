package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zergslaw/boilerplate/internal/app"
)

// SaveCode need for implements app.CodeRepo.
func (repo *Repo) SaveCode(ctx context.Context, userID app.UserID, code string) error {
	return repo.execFunc(func(db *sql.DB) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("start tx: %w", err)
		}

		err = cleanRecoveryCodes(ctx, tx, userID)
		if err != nil {
			return rollback(tx, err)
		}

		const query = `INSERT INTO recovery_code(user_id, code) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, userID, code)
		if err != nil {
			return rollback(tx, err)
		}

		err = createTaskNotification(ctx, tx, userID, app.PassRecovery)
		if err != nil {
			return rollback(tx, err)
		}

		return tx.Commit()
	})
}

// UserIDByCode need for implements app.CodeRepo.
func (repo *Repo) UserIDByCode(ctx context.Context, code string) (userID app.UserID, createAt time.Time, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT user_id, created_at FROM recovery_code WHERE code = $1`

		err = db.QueryRowContext(ctx, query, code).Scan(&userID, &createAt)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return err
		}

		return nil
	})
	return
}

// Code need for implements app.CodeRepo.
func (repo *Repo) Code(ctx context.Context, id app.UserID) (code string, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT code FROM recovery_code WHERE user_id = $1`

		err = db.QueryRowContext(ctx, query, id).Scan(&code)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return err
		}

		return nil
	})
	return
}
