package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockCommandRunner struct {
	ExpectedCommand ExpectedCommand
	T               *testing.T
}

type ExpectedCommand struct {
	Name   string
	Args   []string
	Output []byte
	Error  error
}

func (m *MockCommandRunner) Run(name string, args ...string) ([]byte, error) {
	assert.Equal(m.T, m.ExpectedCommand.Name, name)
	assert.Equal(m.T, m.ExpectedCommand.Args, args)

	return m.ExpectedCommand.Output, m.ExpectedCommand.Error
}
