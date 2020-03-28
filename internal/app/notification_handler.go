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

func (a *Application) execNotification(ctx context.Context, task TaskNotification) error {
	user, err := a.userRepo.UserByID(ctx, task.UserID)
	if err != nil {
		return err
	}

	switch task.Kind {
	case Welcome:
		err = a.sendNotification(Welcome, user.Email, welcomeMsg)
	case ChangeEmail:
		err = a.sendNotification(ChangeEmail, user.Email, changeEmailMsg)
	case PassRecovery:
		err = a.sendRecoveryCode(ctx, user.Email, task)
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

func (a *Application) sendRecoveryCode(ctx context.Context, contact string, task TaskNotification) error {
	code, err := a.codeRepo.Code(ctx, task.UserID)
	if err != nil {
		return err
	}

	return a.sendNotification(task.Kind, contact, code)
}

func (a *Application) sendNotification(kind MessageKind, contact, content string) error {
	return a.notification.Notification(contact, Message{
		Kind:    kind,
		Content: content,
	})
}
