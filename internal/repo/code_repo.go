package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/zergslaw/boilerplate/internal/app"
)

// SaveCode need for implements app.CodeRepo.
func (repo *Repo) SaveCode(ctx context.Context, userID app.UserID, code string) error {
	return repo.db.Tx(ctx, func(tx *sql.Tx) error {
		err := cleanRecoveryCodes(ctx, tx, userID)
		if err != nil {
			return err
		}

		const query = `INSERT INTO recovery_code(user_id, code) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, userID, code)
		if err != nil {
			return err
		}

		err = createTaskNotification(ctx, tx, userID, app.PassRecovery)
		if err != nil {
			return err
		}

		return nil
	})
}

// UserIDByCode need for implements app.CodeRepo.
func (repo *Repo) UserIDByCode(ctx context.Context, code string) (userID app.UserID, createAt time.Time, err error) {
	err = repo.db.Do(func(db *sql.DB) error {
		const query = `SELECT user_id, created_at FROM recovery_code WHERE code = $1`

		err = db.QueryRowContext(ctx, query, code).Scan(&userID, &createAt)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

// Code need for implements app.CodeRepo.
func (repo *Repo) Code(ctx context.Context, id app.UserID) (code string, err error) {
	err = repo.db.Do(func(db *sql.DB) error {
		const query = `SELECT code FROM recovery_code WHERE user_id = $1`

		err = db.QueryRowContext(ctx, query, id).Scan(&code)
		if err != nil {
			return err
		}

		return nil
	})
	return
}
