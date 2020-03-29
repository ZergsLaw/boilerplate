// +build integration

package repo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestUserRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()

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

	user2 := userGenerator()
	user2.ID, err = Repo.CreateUser(ctx, user2)
	assert.Nil(t, err)
	assert.NotZero(t, user2.ID)

	res, err = Repo.UserByUsername(ctx, user2.Username)
	assert.Nil(t, err)
	user2.CreatedAt = res.CreatedAt
	user2.UpdatedAt = res.UpdatedAt
	assert.Equal(t, &user2, res)

	user3 := userGenerator()
	user3.ID, err = Repo.CreateUser(ctx, user3)
	assert.Nil(t, err)
	assert.NotZero(t, user3.ID)

	err = Repo.DeleteUser(ctx, 115)
	assert.Nil(t, err)

	users, total, err := Repo.ListUserByUsername(ctx, "username", app.Page{Limit: 10})
	assert.Nil(t, err)
	user3.CreatedAt = users[0].CreatedAt
	user3.UpdatedAt = users[0].UpdatedAt
	assert.Equal(t, []app.User{user3, user2}, users)
	assert.Equal(t, 2, total)
}
