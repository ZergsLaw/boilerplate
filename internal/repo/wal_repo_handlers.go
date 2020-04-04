package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zergslaw/boilerplate/internal/app"
)

// NotificationTask need for implements app.WAL.
func (repo *Repo) NotificationTask(ctx context.Context) (task *app.TaskNotification, err error) {
	err = repo.execFunc(func(db *sql.DB) error {
		const query = `SELECT id, kind, user_id FROM notifications
		WHERE is_done = false 
		ORDER BY created_at LIMIT 1`

		res := &taskNotificationDBFormat{}
		err = db.QueryRowContext(ctx, query).Scan(&res.ID, &res.Kind, &res.UserID)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return fmt.Errorf("get notification task: %w", err)
		}

		task = res.toAppFormat()
		return nil
	})
	return task, err
}

// DeleteTaskNotification need for implements app.WAL.
func (repo *Repo) DeleteTaskNotification(ctx context.Context, id int) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `UPDATE notifications SET is_done = true, exec_time = now() WHERE id = $1`

		_, err := db.ExecContext(ctx, query, id)

		return err
	})
}
