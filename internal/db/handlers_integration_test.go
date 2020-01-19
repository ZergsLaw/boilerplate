// +build integration

package db_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/users/internal/app"
)

var (
	userGenerator = generatorUser()
	ctx           = context.Background()
	ip            = "192.100.10.4"
	origin        = app.Origin{
		IP:        net.ParseIP(ip),
		UserAgent: "UserAgent",
	}
)

func TestRepoSmoke(t *testing.T) {
	user := userGenerator()

	var err error
	user.ID, err = Repo.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	res, err := Repo.UserByID(ctx, user.ID)
	assert.Nil(t, err)
	user.CreatedAt = res.CreatedAt
	user.UpdatedAt = res.UpdatedAt
	assert.Equal(t, &user, res)

	newUsername := "newUsername"
	err = Repo.UpdateUsername(ctx, user.ID, newUsername)
	assert.Nil(t, err)
	user.Username = newUsername

	newEmail := "newEmail@gmail.com"
	err = Repo.UpdateEmail(ctx, user.ID, newEmail)
	assert.Nil(t, err)
	user.Email = newEmail

	res, err = Repo.UserByEmail(ctx, user.Email)
	assert.Nil(t, err)
	user.UpdatedAt = res.UpdatedAt
	assert.Equal(t, &user, res)

	newPass := []byte(`newPassword`)
	err = Repo.UpdatePassword(ctx, user.ID, newPass)
	assert.Nil(t, err)
	user.PassHash = newPass

	res, err = Repo.UserByUsername(ctx, user.Username)
	assert.Nil(t, err)
	user.UpdatedAt = res.UpdatedAt
	assert.Equal(t, &user, res)

	user2 := userGenerator()
	user2.ID, err = Repo.CreateUser(ctx, user2)
	assert.Nil(t, err)
	assert.NotZero(t, user2.ID)

	user3 := userGenerator()
	user3.ID, err = Repo.CreateUser(ctx, user3)
	assert.Nil(t, err)
	assert.NotZero(t, user3.ID)

	const tokenUser2 = "token2"
	err = Repo.SaveSession(ctx, user2.ID, tokenUser2, origin)
	assert.Nil(t, err)

	expectedSession := &app.Session{
		Origin:  origin,
		TokenID: tokenUser2,
	}

	session, err := Repo.SessionByTokenID(ctx, tokenUser2)
	assert.Nil(t, err)
	expectedSession.ID = session.ID
	if expectedSession.IP.Equal(session.IP) {
		expectedSession.IP = session.IP
	}
	assert.Equal(t, expectedSession, session)

	res, err = Repo.UserByTokenID(ctx, tokenUser2)
	assert.Nil(t, err)
	user2.CreatedAt = res.CreatedAt
	user2.UpdatedAt = res.UpdatedAt
	assert.Equal(t, &user2, res)

	err = Repo.DeleteUser(ctx, 115)
	assert.Nil(t, err)

	users, total, err := Repo.ListUserByUsername(ctx, "username", app.Page{Limit: 10})
	assert.Nil(t, err)
	user3.CreatedAt = users[0].CreatedAt
	user3.UpdatedAt = users[0].UpdatedAt
	assert.Equal(t, []app.User{user3, user2}, users)
	assert.Equal(t, 2, total)

	err = Repo.DeleteSession(ctx, tokenUser2)
	assert.Nil(t, err)

	session, err = Repo.SessionByTokenID(ctx, tokenUser2)
	assert.Nil(t, session)
	assert.Equal(t, app.ErrNotFound, errors.Unwrap(err))
}

func generatorUser() func() app.User {
	x := 0

	return func() app.User {
		x++
		return app.User{
			ID:        app.UserID(x),
			Email:     fmt.Sprintf("email%d@gmail.com", x),
			Username:  fmt.Sprintf("username%d", x),
			PassHash:  []byte("pass"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
}
