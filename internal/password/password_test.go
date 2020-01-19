package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zergslaw/users/internal/password"
)

var (
	pass = "pass"
)

func TestPassword(t *testing.T) {
	t.Parallel()

	passwords := password.New()
	hashPass, err := passwords.Hashing(pass)
	assert.NoError(t, err)
	compare := passwords.Compare(hashPass, []byte(pass))
	assert.Equal(t, true, compare)
}
