// Package auth contains methods for working with authorization tokens,
// their generation and parsing.
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/zergslaw/users/internal/app"
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

const (
	// CookieTokenName name for auth cookie
	CookieTokenName = "__Secure-authKey" //nolint:gosec
)

// Token need for implements app.Auth.
func (t *Auth) Token(expired time.Duration) (app.AuthToken, app.TokenID, error) {
	tokenID, err := t.generatorID()
	if err != nil {
		return "", "", fmt.Errorf("uuid generate: %w", err)
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

	cookie := http.Cookie{
		Name:     CookieTokenName,
		Value:    tokenString,
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return app.AuthToken(cookie.String()), app.TokenID(tokenID), nil
}

// Parse need for implements app.Auth.
func (t *Auth) Parse(authToken app.AuthToken) (app.TokenID, error) {
	tokenString := parseToken(authToken)
	if tokenString == "" {
		return "", app.ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func parseToken(authToken app.AuthToken) string {
	header := http.Header{}
	header.Add("Cookie", string(authToken))
	request := http.Request{Header: header}

	cookieKey, err := request.Cookie(CookieTokenName)
	if err != nil {
		return ""
	}

	return cookieKey.Value
}
