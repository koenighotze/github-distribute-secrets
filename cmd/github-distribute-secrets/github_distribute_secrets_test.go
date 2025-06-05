package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"koenighotze.de/github-distribute-secrets/internal/config"
)

type MockOnePasswordClient struct {
	expectedError error
	calls         int
}

func (m *MockOnePasswordClient) GetSecret(secretPath string) (secret string, err error) {
	m.calls++
	return "something", m.expectedError
}

type mockGithubClient struct {
	calls         int
	expectedError error
}

func (m *mockGithubClient) AddSecretToRepository(key string, secret string, repositoy string) (err error) {
	m.calls++
	return m.expectedError
}

type MockConfigFileReader struct {
	expectedConfig *config.Configuration
	expectedError  error
	calls          int
}

func (m *MockConfigFileReader) ReadConfiguration(path string) (config *config.Configuration, err error) {
	m.calls++
	return m.expectedConfig, m.expectedError
}

func TestApplyConfigurationToRepository(t *testing.T) {
	configMap := config.RepositoryConfiguration{
		"foo": "bar",
	}
	repository := "aname"

	t.Run("should not add a secret to the repo if reading the secret failed", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		onePasswordClient.expectedError = assert.AnError

		_ = applyConfigurationToRepository(configMap, repository, onePasswordClient, githubClient)

		assert.Equal(t, 1, onePasswordClient.calls)
		assert.Equal(t, 0, githubClient.calls)
	})

	t.Run("should return true if the config map is empty", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}

		result := applyConfigurationToRepository(config.RepositoryConfiguration{}, repository, onePasswordClient, githubClient)

		assert.True(t, result)
		assert.Equal(t, 0, onePasswordClient.calls)
		assert.Equal(t, 0, githubClient.calls)
	})

	t.Run("should return true if all secrets where applied successfully", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}

		result := applyConfigurationToRepository(configMap, repository, onePasswordClient, githubClient)

		assert.True(t, result)
		assert.Equal(t, 1, onePasswordClient.calls)
		assert.Equal(t, 1, githubClient.calls)
	})

	t.Run("should return false if at least one secret was not applied successfully", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		githubClient.expectedError = assert.AnError

		result := applyConfigurationToRepository(configMap, repository, onePasswordClient, githubClient)

		assert.False(t, result)
		assert.Equal(t, 1, onePasswordClient.calls)
		assert.Equal(t, 1, githubClient.calls)
	})

	t.Run("should return false if at least one secret could not be read", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		onePasswordClient.expectedError = assert.AnError

		result := applyConfigurationToRepository(configMap, repository, onePasswordClient, githubClient)

		assert.False(t, result)
		assert.Equal(t, 1, onePasswordClient.calls)
		assert.Equal(t, 0, githubClient.calls)
	})

	t.Run("should apply all secrets of the config map to the repository", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configMap := config.RepositoryConfiguration{
			"foo": "bar",
			"faz": "fumm",
		}

		result := applyConfigurationToRepository(configMap, repository, onePasswordClient, githubClient)

		assert.True(t, result)
		assert.Equal(t, 2, onePasswordClient.calls)
		assert.Equal(t, 2, githubClient.calls)
	})
}

func TestApplyConfiguration(t *testing.T) {
	configuration := &config.Configuration{
		RawConfig: map[string]config.RepositoryConfiguration{
			"repo1": {
				"key": "val",
			},
		},
		Repositories: []string{"repo1"},
	}

	t.Run("should return true if no errors occured", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}

		result := applyConfiguration(configuration, onePasswordClient, githubClient)

		assert.True(t, result)
	})

	t.Run("should return false if at least one error occured", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		onePasswordClient.expectedError = assert.AnError

		result := applyConfiguration(configuration, onePasswordClient, githubClient)

		assert.False(t, result)
	})

	t.Run("should return true if repositories are empty", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configuration.Repositories = []string{}

		result := applyConfiguration(configuration, onePasswordClient, githubClient)

		assert.True(t, result)
	})

	t.Run("should apply the configuration to all repositories", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configuration.RawConfig = map[string]config.RepositoryConfiguration{
			"foo": map[string]string{"k": "v"},
			"bar": map[string]string{"k": "v"},
			"baz": map[string]string{"k": "v"},
		}
		configuration.Repositories = []string{"foo", "bar", "baz"}

		_ = applyConfiguration(configuration, onePasswordClient, githubClient)

		assert.Equal(t, len(configuration.Repositories), githubClient.calls)
	})
}

func TestGithubSecretDistribution(t *testing.T) {
	configuration := &config.Configuration{
		RawConfig: map[string]config.RepositoryConfiguration{
			"repo1": {
				"key": "val",
			},
		},
		Repositories: []string{"repo1"},
	}
	t.Run("should read the configuration", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configFileReader := &MockConfigFileReader{
			expectedConfig: configuration,
		}

		githubSecretDistribution(configFileReader, onePasswordClient, githubClient)

		assert.Equal(t, 1, configFileReader.calls)
	})

	t.Run("should apply the configuration", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configFileReader := &MockConfigFileReader{
			expectedConfig: configuration,
		}

		githubSecretDistribution(configFileReader, onePasswordClient, githubClient)

		assert.Equal(t, 1, githubClient.calls)
	})

	t.Run("should panic even if a single application fails", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{
			expectedError: assert.AnError,
		}
		configFileReader := &MockConfigFileReader{
			expectedConfig: configuration,
		}

		assert.Panics(t, func() {
			githubSecretDistribution(configFileReader, onePasswordClient, githubClient)
		})
	})

	t.Run("should panic if reading the config failed", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configFileReader := &MockConfigFileReader{
			expectedError: assert.AnError,
		}

		assert.Panics(t, func() {
			githubSecretDistribution(configFileReader, onePasswordClient, githubClient)
		})
	})
}
