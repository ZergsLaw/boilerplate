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

		id, userID, kind := 0, app.UserID(0), ""
		err = db.QueryRowContext(ctx, query).Scan(&id, &kind, &userID)
		switch {
		case err == sql.ErrNoRows:
			return app.ErrNotFound
		case err != nil:
			return fmt.Errorf("get notification task: %w", err)
		}

		msgKind, err := parseKindNotification(kind)
		if err != nil {
			return err
		}

		task = &app.TaskNotification{
			ID:     id,
			UserID: userID,
			Kind:   msgKind,
		}
		return nil
	})
	return task, err
}

func parseKindNotification(str string) (app.MessageKind, error) {
	switch str {
	case app.Welcome.String():
		return app.Welcome, nil
	case app.ChangeEmail.String():
		return app.ChangeEmail, nil
	case app.PassRecovery.String():
		return app.PassRecovery, nil
	default:
		return 0, app.ErrNotUnknownKindTask
	}
}

// DeleteTaskNotification need for implements app.WAL.
func (repo *Repo) DeleteTaskNotification(ctx context.Context, id int) error {
	return repo.execFunc(func(db *sql.DB) error {
		const query = `UPDATE notifications SET is_done = true, exec_time = now() WHERE id = $1`

		_, err := db.ExecContext(ctx, query, id)

		return err
	})
}
