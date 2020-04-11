package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zergslaw/boilerplate/internal/app"
)

func createTaskNotification(ctx context.Context, tx *sql.Tx, userID app.UserID, kind app.MessageKind) error {
	const queryCreateTask = `INSERT INTO notifications (user_id, kind) VALUES ($1, $2)`

	_, err := tx.ExecContext(ctx, queryCreateTask, userID, kind.String())
	if err != nil {
		return fmt.Errorf("create task notification: %w", err)
	}

	return nil
}

func cleanRecoveryCodes(ctx context.Context, tx *sql.Tx, id app.UserID) error {
	const query = `DELETE FROM recovery_code WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete recovery recoverycode: %w", err)
	}

	return nil
}

func rowsCloseWithError(rows *sql.Rows, err error) error {
	errFromRows := rows.Close()
	if errFromRows != nil {
		err = fmt.Errorf("error: %w, rows close: %s", err, errFromRows)
	}

	return err
}
