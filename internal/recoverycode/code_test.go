package recoverycode_test

import (
	"testing"

	"github.com/zergslaw/boilerplate/internal/recoverycode"
)

func TestCode_Generate(t *testing.T) {
	t.Parallel()

	c := recoverycode.New()

	const countIteration = 1000
	const length = 4
	result := make(map[string]bool)
	for i := 1; i <= countIteration; i++ {
		newCode := c.Generate(length)
		if result[newCode] {
			t.Fatal("the same code was generated")
		}

		result[newCode] = true
	}
}
