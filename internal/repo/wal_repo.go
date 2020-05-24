package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zergslaw/boilerplate/internal/app"
)

// NotificationTask need for implements app.WAL.
func (repo *Repo) NotificationTask(ctx context.Context) (task *app.TaskNotification, err error) {
	err = repo.db.Do(func(db *sqlx.DB) error {
		const query = `SELECT id, kind, email FROM notifications
		WHERE is_done = false 
		ORDER BY created_at LIMIT 1`

		res := &taskNotificationDBFormat{}
		err = db.GetContext(ctx, res, query)
		if err != nil {
			return err
		}

		task = res.toAppFormat()
		return nil
	})
	return
}

// DeleteTaskNotification need for implements app.WAL.
func (repo *Repo) DeleteTaskNotification(ctx context.Context, id int) error {
	return repo.db.Do(func(db *sqlx.DB) error {
		const query = `UPDATE notifications SET is_done = true, exec_time = now() WHERE id = $1`

		_, err := db.ExecContext(ctx, query, id)

		return err
	})
}
