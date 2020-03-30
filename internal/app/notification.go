package app

import (
	"context"
	"fmt"
	"time"
)

type (
	// Notification module for working with alerts for registered users.
	Notification interface {
		// NotificationTask must accept the parameter contact to whom the notification will be sent.
		// At the moment, the guarantee of message delivery lies on this module, it is possible to
		// transfer it to the Application.
		Notification(contact string, msg Message) error
	}
	// Message contains sent info.
	Message struct {
		Kind    MessageKind
		Content string
	}
	// TaskNotification contains information to perform the task of notifying the user.
	TaskNotification struct {
		ID     int
		UserID UserID
		Kind   MessageKind
	}
	// MessageKind selects the type of message to be sent.
	MessageKind int
)

// Message enums.
const (
	Welcome MessageKind = iota + 1
	ChangeEmail
	PassRecovery
)

func (m MessageKind) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case ChangeEmail:
		return "change email"
	case PassRecovery:
		return "password recovery"
	default:
		panic(fmt.Sprintf("unknown kind: %d", m))
	}
}

func wait(ctx context.Context) {
	const timeDelay = time.Second

	select {
	case <-ctx.Done():
	case <-time.After(timeDelay):
	}
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
