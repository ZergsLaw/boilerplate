package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/users/internal/app"
	"github.com/zergslaw/users/internal/auth"
)

var (
	expired                = time.Hour
	appTokenID app.TokenID = "tokenID"

	generateID = func() (string, error) {
		return string(appTokenID), nil
	}
)

func TestToken(t *testing.T) {
	t.Parallel()

	tokenizer := auth.New("super-duper-secret-key", auth.SetIDGenerator(generateID))

	appToken, tokenID, err := tokenizer.Token(expired)
	assert.NoError(t, err)
	assert.NotZero(t, appToken)
	assert.Equal(t, appTokenID, tokenID)

	tokenID, err = tokenizer.Parse(appToken)
	assert.NoError(t, err)
	assert.Equal(t, appTokenID, tokenID)
}
