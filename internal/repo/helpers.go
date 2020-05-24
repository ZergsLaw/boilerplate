package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/zergslaw/boilerplate/internal/app"
)

func createTaskNotification(ctx context.Context, tx *sqlx.Tx, task app.TaskNotification) error {
	const queryCreateTask = `INSERT INTO notifications (email, kind) VALUES (:email, :kind)`
	type args struct {
		Email string `db:"email"`
		Kind  string `db:"kind"`
	}

	_, err := tx.NamedExecContext(ctx, queryCreateTask, args{
		Email: task.Email,
		Kind:  task.Kind.String(),
	})
	if err != nil {
		return fmt.Errorf("create task notification: %w", err)
	}

	return nil
}

func cleanRecoveryCodes(ctx context.Context, tx *sqlx.Tx, email string) error {
	const query = `DELETE FROM recovery_code WHERE email = $1`

	_, err := tx.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("delete recovery recoverycode: %w", err)
	}

	return nil
}
