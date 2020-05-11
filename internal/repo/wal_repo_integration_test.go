// +build integration

package repo_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestWALRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()
	task, err := Repo.NotificationTask(ctx)
	assert.Equal(t, app.ErrNotFound, errors.Unwrap(err))
	assert.Nil(t, task)

	user.ID, err = Repo.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	task, err = Repo.NotificationTask(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, app.Welcome, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	assert.Nil(t, err)

	newEmail := "newEmail@gmail.com"
	err = Repo.UpdateEmail(ctx, user.ID, newEmail)
	assert.Nil(t, err)
	user.Email = newEmail

	task, err = Repo.NotificationTask(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, task.ID)
	assert.Equal(t, app.ChangeEmail, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	assert.Nil(t, err)

	newPass := []byte(`newPassword`)
	err = Repo.UpdatePassword(ctx, user.ID, newPass)
	assert.Nil(t, err)
	user.PassHash = newPass

	const recoveryCode = "123456"
	err = Repo.SaveCode(ctx, user.ID, recoveryCode)
	assert.Nil(t, err)

	task, err = Repo.NotificationTask(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 3, task.ID)
	assert.Equal(t, app.PassRecovery, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	assert.Nil(t, err)
}
