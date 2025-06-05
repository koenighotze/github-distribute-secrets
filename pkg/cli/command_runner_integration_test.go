//go:build integration

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExecIntegration(t *testing.T) {
	t.Run("should return the combined output", func(t *testing.T) {
		output, error := defaultExec("echo", "foo")

		assert.NoError(t, error)
		assert.Equal(t, "foo\n", string(output))
	})

	t.Run("should return the error of the command", func(t *testing.T) {
		_, error := defaultExec("nosuchcommand")

		assert.ErrorContains(t, error, "not found")
	})
}
