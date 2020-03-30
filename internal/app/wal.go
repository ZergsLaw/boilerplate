package app

import (
	"context"
	"errors"
)

type (
	// WALApplication a provider to run tasks.
	WALApplication interface {
		// StartWALNotification starts the task of notifying users.
		StartWALNotification(ctx context.Context) error
	}
	// WAL module returning tasks and also closing them.
	WAL interface {
		// NotificationTask returns the earliest task that has not been completed.
		// Errors: ErrNotFound, unknown.
		NotificationTask(ctx context.Context) (task *TaskNotification, err error)
		// DeleteTaskNotification removes the task performed.
		// Errors: unknown.
		DeleteTaskNotification(ctx context.Context, id int) error
	}
)

// StartWALNotification for implemented WALApplication.
func (a *Application) StartWALNotification(ctx context.Context) error {
	for ctx.Err() == nil {
		task, err := a.wal.NotificationTask(ctx)
		switch {
		case err == nil:
			err := a.execNotification(ctx, *task)
			if err != nil {
				return err
			}
		case errors.Is(err, ErrNotFound):
			wait(ctx)
		default:
			return err
		}
	}

	return ctx.Err()
}
