// +build integration

package repo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestCodeRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()
	user.ID, err = Repo.CreateUser(ctx, user, app.TaskNotification{
		Email: user.Email,
		Kind:  app.Welcome,
	})
	require.Nil(t, err)
	require.NotZero(t, user.ID)

	const recoveryCode = "123456"
	err = Repo.SaveCode(ctx, user.Email, recoveryCode, app.TaskNotification{
		Email: user.Email,
		Kind:  app.PassRecovery,
	})
	require.Nil(t, err)

	codeInfo, err := Repo.Code(ctx, user.Email)
	require.Nil(t, err)
	expected := &app.CodeInfo{
		Code:      recoveryCode,
		Email:     user.Email,
		CreatedAt: codeInfo.CreatedAt,
	}
	require.Equal(t, expected, codeInfo)
}
