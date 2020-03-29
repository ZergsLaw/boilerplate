// +build integration

package repo_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zergslaw/boilerplate/internal/app"
)

func TestCodeRepoSmoke(t *testing.T) {
	err := truncate()
	require.NoError(t, err)

	user := userGenerator()
	user.ID, err = Repo.CreateUser(ctx, user)
	assert.Nil(t, err)
	assert.NotZero(t, user.ID)

	const recoveryCode = "123456"
	userID, createdAt, err := Repo.UserIDByCode(ctx, recoveryCode)
	assert.Zero(t, userID)
	assert.Zero(t, createdAt)
	assert.Equal(t, app.ErrNotFound, errors.Unwrap(err))

	err = Repo.SaveCode(ctx, user.ID, recoveryCode)
	assert.Nil(t, err)

	code, err := Repo.Code(ctx, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, recoveryCode, code)

	userID, createdAt, err = Repo.UserIDByCode(ctx, recoveryCode)
	assert.Nil(t, err)
	assert.NotZero(t, createdAt)
	assert.Equal(t, user.ID, userID)
}
