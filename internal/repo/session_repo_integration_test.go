// +build integration

package repo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestSessionRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()

	user.ID, err = Repo.CreateUser(ctx, user, app.TaskNotification{
		Email: user.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user.ID)

	const tokenUser = "token"
	err = Repo.SaveSession(ctx, user.ID, tokenUser, origin)
	require.Nil(t, err)

	expectedSession := &app.Session{
		Origin:  origin,
		TokenID: tokenUser,
	}

	session, err := Repo.SessionByTokenID(ctx, tokenUser)
	require.Nil(t, err)
	expectedSession.ID = session.ID
	if expectedSession.IP.Equal(session.IP) {
		expectedSession.IP = session.IP
	}
	require.Equal(t, expectedSession, session)

	userFromDB, err := Repo.UserByTokenID(ctx, tokenUser)
	require.Nil(t, err)
	user.CreatedAt = userFromDB.CreatedAt
	user.UpdatedAt = userFromDB.UpdatedAt
	require.Equal(t, user, *userFromDB)

	err = Repo.DeleteSession(ctx, tokenUser)
	require.Nil(t, err)
}
