package app

import (
	"context"
	"errors"
	"time"
)

func wait(ctx context.Context) {
	const timeDelay = time.Second

	select {
	case <-ctx.Done():
	case <-time.After(timeDelay):
	}
}

func (a *app) StartWALNotification(ctx context.Context) error {
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

func (a *app) execNotification(ctx context.Context, task TaskNotification) (err error) {
	switch task.Kind {
	case Welcome:
		err = a.welcome(task.Email)
	case ChangeEmail:
		err = a.changeEmail(task.Email)
	default:
		err = ErrNotUnknownKindTask
	}
	if err != nil {
		return err
	}

	return a.wal.DeleteTaskNotification(ctx, task.ID)
}

const (
	welcomeMsg     = `Welcome`
	changeEmailMsg = `Change email successful`
)

func (a *app) welcome(contact string) error {
	return a.notification.Notification(contact, Message{
		Kind:    Welcome,
		Content: welcomeMsg,
	})
}

func (a *app) changeEmail(contact string) error {
	return a.notification.Notification(contact, Message{
		Kind:    ChangeEmail,
		Content: changeEmailMsg,
	})
}
