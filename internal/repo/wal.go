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
		const query = `SELECT notifications.id, notifications.kind, users.email FROM
        notifications LEFT JOIN users ON notifications.user_id = users.id 
		WHERE notifications.is_done = false 
		ORDER BY notifications.created_at LIMIT 1`

		id, email, kind := 0, "", ""
		err = db.QueryRowContext(ctx, query).Scan(&id, &kind, &email)
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
			ID:    id,
			Email: email,
			Kind:  msgKind,
		}
		return nil
	})
	return task, err
}

func parseKindNotification(str string) (app.Message, error) {
	switch str {
	case app.Welcome.String():
		return app.Welcome, nil
	case app.ChangeEmail.String():
		return app.ChangeEmail, nil
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
