package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zergslaw/users/internal/app"
)

const (
	cookieTokenName = "__Secure-authKey" // nolint:gosec
	authTimeout     = 250 * time.Millisecond
)

func (svc *service) cookieKeyAuth(raw string) (*app.AuthUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), authTimeout)
	defer cancel()
	profile, err := svc.app.UserByAuthToken(ctx, parseToken(raw))
	switch {
	case err != nil:
		return nil, fmt.Errorf("userByAuthToken: %w", err)
	default:
		return profile, nil
	}
}

func parseToken(raw string) app.AuthToken {
	header := http.Header{}
	header.Add("Cookie", raw)
	request := http.Request{Header: header}
	cookieKey, err := request.Cookie(cookieTokenName)
	if err != nil {
		return ""
	}

	return app.AuthToken(cookieKey.Value)
}

func generateCookie(token app.AuthToken) *http.Cookie {
	cookie := &http.Cookie{
		Name:     cookieTokenName,
		Value:    string(token),
		Secure:   true,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}
