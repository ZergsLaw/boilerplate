// +build integration

package repo_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestWALRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()
	task, err := Repo.NotificationTask(ctx)
	require.Equal(t, app.ErrNotFound, errors.Unwrap(err))
	require.Nil(t, task)

	user.ID, err = Repo.CreateUser(ctx, user, app.TaskNotification{
		Email: user.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user.ID)

	task, err = Repo.NotificationTask(ctx)
	require.Nil(t, err)
	require.Equal(t, 1, task.ID)
	require.Equal(t, app.Welcome, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	require.Nil(t, err)

	newEmail := "newEmail@gmail.com"
	err = Repo.UpdateEmail(ctx, user.ID, newEmail, app.TaskNotification{
		Email: newEmail,
		Kind:  app.ChangeEmail,
	})
	require.Nil(t, err)
	user.Email = newEmail

	task, err = Repo.NotificationTask(ctx)
	require.Nil(t, err)
	require.Equal(t, 2, task.ID)
	require.Equal(t, app.ChangeEmail, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	require.Nil(t, err)

	newPass := []byte(`newPassword`)
	err = Repo.UpdatePassword(ctx, user.ID, newPass)
	require.Nil(t, err)
	user.PassHash = newPass

	const recoveryCode = "123456"
	err = Repo.SaveCode(ctx, user.Email, recoveryCode, app.TaskNotification{
		Email: user.Email,
		Kind:  app.PassRecovery,
	})
	require.Nil(t, err)

	task, err = Repo.NotificationTask(ctx)
	require.Nil(t, err)
	require.Equal(t, 3, task.ID)
	require.Equal(t, app.PassRecovery, task.Kind)

	err = Repo.DeleteTaskNotification(ctx, task.ID)
	require.Nil(t, err)
}
