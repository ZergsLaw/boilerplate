// Package recoverycode contains an implementation to generated random lines of code.
package recoverycode

import (
	"math/rand"
	"time"

	"github.com/zergslaw/boilerplate/internal/app"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type (
	// Option for building code struct.
	Option func(*code)

	code struct {
		randInt func(max int) int
	}
)

// Generate need for implemented app.Code.
func (c *code) Generate(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[c.randInt(len(charset))]
	}
	return string(b)
}

// New creates a new instance of the app.Code object.
func New(option ...Option) app.Code {
	c := defaultConfig()

	for i := range option {
		option[i](c)
	}

	return c
}

// RandInt option for sets generator random int.
func RandInt(randInt func(max int) int) Option {
	return func(c *code) {
		c.randInt = randInt
	}
}

func defaultConfig() *code {
	return &code{
		randInt: rand.New(rand.NewSource(time.Now().UnixNano())).Intn,
	}
}
