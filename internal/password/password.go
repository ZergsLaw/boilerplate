// Package password contains methods for hashing and comparing passwords.
package password

import (
	"github.com/zergslaw/users/internal/app"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Password is an implements app.Password.
	// Responsible for working passwords, hashing and compare.
	Password struct {
		cost int
	}
	// Option for building Password struct.
	Option func(*Password)
)

// Cost option for sets hashing cost.
func Cost(cost int) Option {
	return func(password *Password) {
		password.cost = cost
	}
}

// New creates and returns new app.Password.
func New(options ...Option) app.Password {
	p := &Password{cost: bcrypt.DefaultCost}

	for i := range options {
		options[i](p)
	}

	return p
}

// Hashing need for implements app.Password.
func (p *Password) Hashing(password string) ([]byte, error) {
	return Hashing(password, p.cost)
}

// Compare need for implements app.Password.
func (p *Password) Compare(hashedPassword []byte, password []byte) bool {
	err := Compare(hashedPassword, password)
	return err == nil
}

// Hashing hashes the password according to the specified cost.
func Hashing(password string, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

// Compare matched passwords.
func Compare(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
