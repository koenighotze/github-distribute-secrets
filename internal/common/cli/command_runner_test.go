package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandRunner(t *testing.T) {
	t.Run("should return a valid command runner", func(t *testing.T) {
		runner := NewCommandRunner()

		assert.NotNil(t, runner)
		_, ok := runner.(CliCommandRunner)
		assert.True(t, ok, "Expected NewCommandRunner to return a CliCommandRunner")
	})

	type mockExec struct {
		called        bool
		exectedReturn []byte
		exectedError  error
		exec          func(name string, arg ...string) ([]byte, error)
	}

	defaultMockExec := mockExec{}
	defaultMockExec.called = false
	defaultMockExec.exec = func(name string, arg ...string) ([]byte, error) {
		defaultMockExec.called = true
		return defaultMockExec.exectedReturn, defaultMockExec.exectedError
	}

	t.Run("should return the output of the executor", func(t *testing.T) {
		defaultMockExec.exectedReturn = []byte("thereturn")
		runner := CliCommandRunner{
			exec: defaultMockExec.exec,
		}

		result, err := runner.Run("foo", "bar")

		assert.Equal(t, defaultMockExec.exectedReturn, result)
		assert.Nil(t, err)
	})
	t.Run("should return the error of the executor", func(t *testing.T) {
		defaultMockExec.exectedReturn = nil
		defaultMockExec.exectedError = assert.AnError
		runner := CliCommandRunner{
			exec: defaultMockExec.exec,
		}

		result, err := runner.Run("foo", "bar")

		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
	})
}
