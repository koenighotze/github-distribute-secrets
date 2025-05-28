package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockCommandRunner(t *testing.T) {
	t.Run("should return expected output for expected commands", func(t *testing.T) {
		expectedOutput := []byte("test output")
		mockRunner := &MockCommandRunner{
			ExpectedCommand: ExpectedCommand{
				Name:   "echo",
				Args:   []string{"hello", "world"},
				Output: expectedOutput,
				Error:  nil,
			},
			T: t,
		}

		output, err := mockRunner.Run("echo", "hello", "world")

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	})

	t.Run("should return expected error for commands configured to fail", func(t *testing.T) {
		expectedError := assert.AnError
		mockRunner := &MockCommandRunner{
			ExpectedCommand: ExpectedCommand{
				Name:   "failing-command",
				Args:   nil,
				Output: nil,
				Error:  expectedError,
			},
			T: t,
		}

		output, err := mockRunner.Run("failing-command")

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, output)
	})
}

func TestNewCommandRunner(t *testing.T) {
	t.Run("should return a valid command runner", func(t *testing.T) {
		runner := NewCommandRunner()

		assert.NotNil(t, runner)
		_, ok := runner.(CliCommandRunner)
		assert.True(t, ok, "Expected NewCommandRunner to return a CliCommandRunner")
	})
}
