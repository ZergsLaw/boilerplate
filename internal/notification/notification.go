// Package notification service communication module.
// Communication takes place via rabbit.
package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
	"github.com/zergslaw/boilerplate/internal/app"
)

type (
	// for convention testing.
	// implemented *amqp.Channel.
	channel interface {
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}

	client struct {
		channel     channel
		exchange    string
		key         string
		appID       string
		mandatory   bool
		immediate   bool
		generatorID func() (string, error)
	}

	notification struct {
		Contact string `json:"contact"`
	}
)

// Config for connections rabbit.
type Config struct {
	User string
	Pass string
	Host string
	Port int
}

// Connect creates a connection to rabbit.
// IMPORTANT: does not declare a queue or anything else, the calling code must do so.
func Connect(cfg Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.User, cfg.Pass, cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}

	return conn, nil
}

// New creates a new instance of the app.NotificationTask object.
// Accepts the parameter interface type object, which is basically implemented *amqp.Channel.
// Accepts the interface for convenient testing.
func New(ch channel, opt ...Option) app.Notification {
	c := defaultClient(ch)

	for i := range opt {
		opt[i](c)
	}

	return c
}

func defaultClient(ch channel) *client {
	return &client{
		channel:     ch,
		exchange:    "",
		key:         defaultKey,
		appID:       defaultAppID,
		mandatory:   false,
		immediate:   false,
		generatorID: generateID,
	}
}

// for convenient testing.
func generateID() (string, error) {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return tokenID.String(), nil
}

const (
	defaultKey   = `notification`
	defaultAppID = `boilerplate`
)

// NotificationTask need for implemented app.NotificationTask.
func (c *client) Notification(contact string, msg app.Message) error {
	n := notification{
		Contact: contact,
	}

	js, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("json marshal notification: %w", err)
	}

	reqID, err := c.generatorID()
	if err != nil {
		return fmt.Errorf("generate req id: %w", err)
	}

	pub := amqp.Publishing{
		ContentType: http.DetectContentType(js),
		MessageId:   reqID,
		Timestamp:   time.Now(),
		Type:        msg.String(),
		AppId:       c.appID,
		Body:        js,
	}

	return c.channel.Publish(c.exchange, c.key, c.mandatory, c.immediate, pub)
}
