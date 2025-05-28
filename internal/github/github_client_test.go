package github

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"koenighotze.de/github-distribute-secrets/internal/common/cli"
)

const (
	testSecretKey   = "TEST_KEY"
	testSecretValue = "test-secret"
	testRepoName    = "test-repo"
)

func TestNewClient(t *testing.T) {
	t.Run("should return a client", func(t *testing.T) {
		result := NewClient()

		_, ok := result.runner.(cli.CliCommandRunner)
		assert.True(t, ok, "Expected runner to be of type cli.CommandRunner")
	})
}

func createMockCommandRunner(t *testing.T, output []byte, err error) cli.CommandRunner {
	return &cli.MockCommandRunner{
		ExpectedCommand: cli.ExpectedCommand{
			Name:   "gh",
			Args:   []string{"secret", "set", testSecretKey, "--body", testSecretValue, "--repo", testRepoName},
			Output: output,
			Error:  err,
		},
		T: t,
	}
}

func TestCreateMockCommandRunner(t *testing.T) {
	t.Run("should create a mock runner with expected configuration", func(t *testing.T) {
		expectedOutput := []byte("test output")
		expectedError := assert.AnError

		mockRunner := createMockCommandRunner(t, expectedOutput, expectedError)

		mockRunnerConcrete, ok := mockRunner.(*cli.MockCommandRunner)
		assert.True(t, ok, "Expected createMockCommandRunner to return a *cli.MockCommandRunner")

		assert.Equal(t, "gh", mockRunnerConcrete.ExpectedCommand.Name,
			"Expected command name to be 'gh'")
		assert.Equal(t, expectedOutput, mockRunnerConcrete.ExpectedCommand.Output,
			"Expected output to match")
		assert.Equal(t, expectedError, mockRunnerConcrete.ExpectedCommand.Error,
			"Expected error to match")
	})
}

func TestAddSecretToRepository(t *testing.T) {
	t.Run("should add a secret to the repository successfully", func(t *testing.T) {
		mockRunner := createMockCommandRunner(t, []byte("Secret added successfully"), nil)
		client := GithubClient{
			runner: mockRunner,
		}

		err := client.AddSecretToRepository(testSecretKey, testSecretValue, testRepoName)

		assert.NoError(t, err, "Expected no error when adding secret")
	})

	t.Run("should return an error if adding the secret fails", func(t *testing.T) {
		mockError := assert.AnError
		mockRunner := createMockCommandRunner(t, []byte("Error adding secret"), mockError)
		client := GithubClient{
			runner: mockRunner,
		}

		err := client.AddSecretToRepository(testSecretKey, testSecretValue, testRepoName)

		assert.Error(t, err, "Expected error when GitHub CLI fails")
		assert.Equal(t, mockError, err, "Expected the error from the command runner")
	})
}
