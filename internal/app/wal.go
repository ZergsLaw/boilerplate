package app

import "context"

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
