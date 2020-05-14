// Package notification service communication module.
package notification

import (
	"fmt"

	"github.com/matcornic/hermes/v2"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/zergslaw/boilerplate/internal/app"
)

type (
	client struct {
		emailClient *sendgrid.Client
		from        string
		hermes      *hermes.Hermes
	}

	notification struct {
		Contact string
		Content string
	}
)

// Connect creates a connection to rabbit.
func Connect(apiKey string) (*sendgrid.Client, error) {
	client := sendgrid.NewSendClient(apiKey)
	return client, nil
}

// New creates a new instance of the app.NotificationTask object.
func New(emailClient *sendgrid.Client, from string) app.Notification {
	return &client{
		emailClient: emailClient,
		from:        from,
		hermes: &hermes.Hermes{
			Theme:         nil,
			TextDirection: "",
			Product: hermes.Product{
				Name:        "Boilerplate",
				Link:        "https://example-hermes.com/",
				Logo:        "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
				Copyright:   "copyright",
				TroubleText: "trouble text",
			},
			DisableCSSInlining: false,
		},
	}
}

const (
	fromName = `boilerplate`
)

// NotificationTask need for implemented app.NotificationTask.
func (c *client) Notification(contact string, msg app.Message) error {
	n := notification{
		Contact: contact,
		Content: msg.Content,
	}

	from := mail.NewEmail(fromName, c.from)
	to := mail.NewEmail(n.Contact, n.Contact)

	email := hermes.Email{
		Body: hermes.Body{
			Name:   subjectByKind(msg.Kind),
			Intros: []string{n.Content},
		},
	}

	htmlContent, err := c.hermes.GenerateHTML(email)
	if err != nil {
		return fmt.Errorf("generated html: %w", err)
	}

	message := mail.NewSingleEmail(from, subjectByKind(msg.Kind), to, "", htmlContent)

	_, err = c.emailClient.Send(message)
	if err != nil {
		return fmt.Errorf("email send: %w", err)
	}

	return nil
}

func subjectByKind(kind app.MessageKind) string {
	switch kind {
	case app.Welcome:
		return "Welcome to boilerplate."
	case app.ChangeEmail:
		return "You have changed your mail."
	case app.PassRecovery:
		return "Recovery password."
	default:
		panic(fmt.Sprintf("unknown kind %s", kind))
	}
}
