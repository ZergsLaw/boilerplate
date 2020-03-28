package app

import "fmt"

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
