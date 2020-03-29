// +build integration

package repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestSessionRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()

	user.ID, err = Repo.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	const tokenUser = "token"
	err = Repo.SaveSession(ctx, user.ID, tokenUser, origin)
	assert.Nil(t, err)

	expectedSession := &app.Session{
		Origin:  origin,
		TokenID: tokenUser,
	}

	session, err := Repo.SessionByTokenID(ctx, tokenUser)
	assert.Nil(t, err)
	expectedSession.ID = session.ID
	if expectedSession.IP.Equal(session.IP) {
		expectedSession.IP = session.IP
	}
	assert.Equal(t, expectedSession, session)

	userFromDB, err := Repo.UserByTokenID(ctx, tokenUser)
	assert.Nil(t, err)
	user.CreatedAt = userFromDB.CreatedAt
	user.UpdatedAt = userFromDB.UpdatedAt
	assert.Equal(t, user, *userFromDB)

	err = Repo.DeleteSession(ctx, tokenUser)
	assert.Nil(t, err)
}
