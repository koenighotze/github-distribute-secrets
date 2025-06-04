package github

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

func TestWithDryRun(t *testing.T) {
	t.Run("should return the dry run client", func(t *testing.T) {
		result := withDryRun()

		_, ok := result.(*dryRunGithubClient)
		assert.True(t, ok, "Expected result to be of type *dryRunGithubClient")
	})
}

func TestDryRunAddSecretToRepository(t *testing.T) {
	t.Run("should return nil if the repository exists", func(t *testing.T) {
		mockRunner := createDryRunMockCommandRunner(t, []byte("Repository exists"), nil)
		client := dryRunGithubClient{
			runner: mockRunner,
		}

		err := client.AddSecretToRepository(testSecretKey, testSecretValue, testRepoName)

		assert.NoError(t, err, "Expected no error in dry run mode for existing repository")
	})

	t.Run("should return an error if the repository does not exist", func(t *testing.T) {
		mockError := fmt.Errorf("repository not found")
		mockRunner := createDryRunMockCommandRunner(t, []byte("Error: repository not found"), mockError)
		client := dryRunGithubClient{
			runner: mockRunner,
		}

		err := client.AddSecretToRepository(testSecretKey, testSecretValue, testRepoName)

		assert.Error(t, err, "Expected error when repository doesn't exist")
		assert.Contains(t, err.Error(), "repository test-repo does not seem to exist")
	})
}

func createDryRunMockCommandRunner(t *testing.T, output []byte, err error) cli.CommandRunner {
	return &cli.MockCommandRunner{
		ExpectedCommand: cli.ExpectedCommand{
			Name:   "gh",
			Args:   []string{"repo", "view", testRepoName},
			Output: output,
			Error:  err,
		},
		T: t,
	}
}
