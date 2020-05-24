// +build integration

package repo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestUserRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()

	user.ID, err = Repo.CreateUser(ctx, user, app.TaskNotification{
		Email: user.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user.ID)

	res, err := Repo.UserByID(ctx, user.ID)
	require.Nil(t, err)
	user.CreatedAt = res.CreatedAt
	user.UpdatedAt = res.UpdatedAt
	require.Equal(t, &user, res)

	newUsername := "newUsername"
	err = Repo.UpdateUsername(ctx, user.ID, newUsername)
	require.Nil(t, err)
	user.Name = newUsername

	newEmail := "newEmail@gmail.com"
	err = Repo.UpdateEmail(ctx, user.ID, newEmail, app.TaskNotification{
		Email: newEmail,
		Kind:  app.ChangeEmail,
	})
	require.Nil(t, err)
	user.Email = newEmail

	res, err = Repo.UserByEmail(ctx, user.Email)
	require.Nil(t, err)
	user.UpdatedAt = res.UpdatedAt
	require.Equal(t, &user, res)

	newPass := []byte(`newPassword`)
	err = Repo.UpdatePassword(ctx, user.ID, newPass)
	require.Nil(t, err)
	user.PassHash = newPass

	user2 := userGenerator()
	user2.ID, err = Repo.CreateUser(ctx, user2, app.TaskNotification{
		Email: user2.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user2.ID)

	res, err = Repo.UserByUsername(ctx, user2.Name)
	require.Nil(t, err)
	user2.CreatedAt = res.CreatedAt
	user2.UpdatedAt = res.UpdatedAt
	require.Equal(t, &user2, res)

	user3 := userGenerator()
	user3.ID, err = Repo.CreateUser(ctx, user3, app.TaskNotification{
		Email: user3.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user3.ID)

	err = Repo.DeleteUser(ctx, 115)
	require.Nil(t, err)

	users, total, err := Repo.ListUserByUsername(ctx, "username", app.Page{Limit: 10})
	require.Nil(t, err)
	user3.CreatedAt = users[0].CreatedAt
	user3.UpdatedAt = users[0].UpdatedAt
	require.Equal(t, []app.User{user3, user2}, users)
	require.Equal(t, 2, total)
}
