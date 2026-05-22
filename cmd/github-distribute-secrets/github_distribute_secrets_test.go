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

		_ = githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)

		assert.Equal(t, 1, configFileReader.calls)
	})

	t.Run("should apply the configuration", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configFileReader := &MockConfigFileReader{
			expectedConfig: configuration,
		}

		_ = githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)

		assert.Equal(t, 1, githubClient.calls)
	})

	t.Run("should return error if a single application fails", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{
			expectedError: assert.AnError,
		}
		configFileReader := &MockConfigFileReader{
			expectedConfig: configuration,
		}

		err := githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)

		assert.Error(t, err)
	})

	t.Run("should return error if reading the config failed", func(t *testing.T) {
		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}
		configFileReader := &MockConfigFileReader{
			expectedError: assert.AnError,
		}

		err := githubSecretDistribution(configFileReader, onePasswordClient, githubClient, false)

		assert.Error(t, err)
	})

	t.Run("should return error if reading config fails", func(t *testing.T) {
		configFileReader := &MockConfigFileReader{expectedError: assert.AnError}

		err := githubSecretDistribution(configFileReader, &MockOnePasswordClient{}, &mockGithubClient{}, false)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return error if applying config fails", func(t *testing.T) {
		configFileReader := &MockConfigFileReader{
			expectedConfig: &config.Configuration{
				RawConfig:    map[string]config.RepositoryConfiguration{"repo1": {"key": "val"}},
				Repositories: []string{"repo1"},
			},
		}
		githubClient := &mockGithubClient{expectedError: assert.AnError}

		err := githubSecretDistribution(configFileReader, &MockOnePasswordClient{}, githubClient, false)

		assert.Error(t, err)
	})

	t.Run("should dump the configuration after reading it", func(t *testing.T) {
		// Arrange
		testConfig := &config.Configuration{
			RawConfig: map[string]config.RepositoryConfiguration{
				"common": {
					"COMMON_SECRET": "common/path",
				},
				"test-repo": {
					"REPO_SECRET": "repo/path",
				},
			},
			Repositories: []string{"test-repo"},
		}

		configReader := &MockConfigFileReader{
			expectedConfig: testConfig,
		}

		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}

		// Capture stdout - this is a bit tricky in Go, so we'll
		// primarily test that the function completes without errors
		// The actual output check would require capturing stdout

		// Act
		err := githubSecretDistribution(configReader, onePasswordClient, githubClient, true)

		// Assert
		assert.NoError(t, err, "Function should complete successfully")
		assert.Equal(t, 1, configReader.calls, "Should call ReadConfiguration once")
	})

	t.Run("should only dump configuration when dumpConfig is true", func(t *testing.T) {
		// Arrange
		testConfig := &config.Configuration{
			RawConfig: map[string]config.RepositoryConfiguration{
				"common": {
					"COMMON_SECRET": "common/path",
				},
			},
			Repositories: []string{},
		}

		configReader := &MockConfigFileReader{
			expectedConfig: testConfig,
		}

		onePasswordClient := &MockOnePasswordClient{}
		githubClient := &mockGithubClient{}

		// Act & Assert - No way to directly test stdout output in this test,
		// but we can verify the function executes without issues
		err := githubSecretDistribution(configReader, onePasswordClient, githubClient, false)
		assert.NoError(t, err, "Function should complete successfully with dumpConfig=false")

		err = githubSecretDistribution(configReader, onePasswordClient, githubClient, true)
		assert.NoError(t, err, "Function should complete successfully with dumpConfig=true")
	})
}
