package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zergslaw/boilerplate/internal/app"
)

// SaveCode need for implements app.CodeRepo.
func (repo *Repo) SaveCode(ctx context.Context, email, code string, task app.TaskNotification) error {
	return repo.db.Tx(ctx, func(tx *sqlx.Tx) error {
		err := cleanRecoveryCodes(ctx, tx, email)
		if err != nil {
			return err
		}

		const query = `INSERT INTO recovery_code(email, code) VALUES (:email, :code)`
		type args struct {
			Email string `db:"email"`
			Code  string `db:"code"`
		}

		_, err = tx.NamedExecContext(ctx, query, args{
			Email: email,
			Code:  code,
		})
		if err != nil {
			return fmt.Errorf("insert code: %w", err)
		}

		return createTaskNotification(ctx, tx, task)
	})
}

// Code need for implements app.CodeRepo.
func (repo *Repo) Code(ctx context.Context, email string) (codeInfo *app.CodeInfo, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT * FROM recovery_code WHERE email = $1`

		c := &codeInfoDBFormat{}
		err = db.GetContext(ctx, c, query, email)
		if err != nil {
			return err
		}

		codeInfo = c.toAppFormat()
		return nil
	})
	return
}
