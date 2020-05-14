// Package auth contains methods for working with authorization tokens,
// their generation and parsing.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/zergslaw/boilerplate/internal/app"
)

type (
	// Auth is an implements app.Auth.
	// Responsible for working with authorization tokens, be it cookies or jwt.
	Auth struct {
		jwtKey      []byte
		generatorID func() (string, error)
	}
	// Option for building auth struct.
	Option func(*Auth)
)

// Errors.
var (
	ErrValidateAlg = errors.New("unexpected signing method")
)

// New creates and returns new app.Auth.
func New(jwtKey string, options ...Option) app.Auth {
	t := &Auth{jwtKey: []byte(jwtKey), generatorID: generateID}

	for i := range options {
		options[i](t)
	}

	return t
}

// for convenient testing.
func generateID() (string, error) {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return tokenID.String(), nil
}

// SetIDGenerator sets func for create tokenID.
func SetIDGenerator(generatorID func() (string, error)) Option {
	return func(token *Auth) {
		token.generatorID = generatorID
	}
}

// Token need for implements app.Auth.
func (t *Auth) Token(expired time.Duration) (app.AuthToken, app.TokenID, error) {
	tokenID, err := t.generatorID()
	if err != nil {
		return "", "", fmt.Errorf("uuid generated: %w", err)
	}

	claims := &jwt.StandardClaims{
		Subject:   tokenID,
		ExpiresAt: time.Now().Add(expired).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtKey)
	if err != nil {
		return "", "", err
	}

	return app.AuthToken(tokenString), app.TokenID(tokenID), nil
}

// Parse need for implements app.Auth.
func (t *Auth) Parse(authToken app.AuthToken) (app.TokenID, error) {
	token, err := jwt.ParseWithClaims(string(authToken), &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrValidateAlg
		}
		return t.jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", app.ErrInvalidToken
	}

	claims := token.Claims.(*jwt.StandardClaims)
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", app.ErrExpiredToken
	}

	return app.TokenID(claims.Subject), nil
}
