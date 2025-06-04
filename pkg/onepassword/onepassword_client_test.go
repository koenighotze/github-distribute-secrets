package onepassword

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

const (
	testSecretPath = "somepath"
)

func createMockOnePasswordCommandRunner(t *testing.T, output []byte, err error) cli.CommandRunner {
	return &cli.MockCommandRunner{
		ExpectedCommand: cli.ExpectedCommand{
			Name:   "op",
			Args:   []string{"read", testSecretPath},
			Output: output,
			Error:  err,
		},
		T: t,
	}
}

func TestNewClient(t *testing.T) {
	t.Run("should return the caching client", func(t *testing.T) {
		result := NewClient()

		_, ok := result.(*cachedClient)

		assert.True(t, ok, "Expected result to be of type cachedClient")
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("should return the secret from onepassword", func(t *testing.T) {
		client := cliClient{
			runner: createMockOnePasswordCommandRunner(t, []byte("supersecret\n"), nil),
		}

		result, err := client.GetSecret(testSecretPath)

		assert.Nil(t, err)
		assert.Equal(t, "supersecret", result)
	})

	t.Run("should trim whitespaces from the secret", func(t *testing.T) {
		client := cliClient{
			runner: createMockOnePasswordCommandRunner(t, []byte("   supersecret   \n"), nil),
		}

		result, err := client.GetSecret(testSecretPath)

		assert.Nil(t, err)
		assert.Equal(t, "supersecret", result)
	})

	t.Run("should return the error if reading from onepassword fails", func(t *testing.T) {
		expectedError := errors.New("bumm")
		client := cliClient{
			runner: createMockOnePasswordCommandRunner(t, nil, expectedError),
		}

		_, err := client.GetSecret(testSecretPath)

		assert.ErrorContains(t, err, "bumm")
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestGetSecretWithCache(t *testing.T) {
	prepareClient := func(path string, output []byte, err error, cache *cacheEntry) cachedClient {
		mockRunner := createMockOnePasswordCommandRunner(t, output, err)
		cliClient := cachedClient{
			Cache: make(secretCacheType),
			Op: &cliClient{
				runner: mockRunner,
			},
		}
		if nil != cache {
			cliClient.Cache[path] = *cache
		}

		return cliClient
	}

	t.Run("should return the uncached value if uncached", func(t *testing.T) {
		cliClient := prepareClient(testSecretPath, []byte("UncachedOutput"), nil, nil)

		result, _ := cliClient.GetSecret(testSecretPath)

		assert.Equal(t, "UncachedOutput", result)
	})

	t.Run("should return the cached value if cached", func(t *testing.T) {
		cliClient := prepareClient(testSecretPath, []byte("UncachedOutput"), nil, &cacheEntry{Value: "cached", Err: nil})

		result, _ := cliClient.GetSecret(testSecretPath)

		assert.Equal(t, "cached", result)
	})

	t.Run("should return the uncached error if uncached", func(t *testing.T) {
		cliClient := prepareClient(testSecretPath, nil, assert.AnError, nil)

		_, err := cliClient.GetSecret(testSecretPath)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should return the cached error if an error occured", func(t *testing.T) {
		expectedError := errors.New("cachederror")
		cliClient := prepareClient(testSecretPath, nil, assert.AnError, &cacheEntry{Value: "cached", Err: expectedError})

		_, err := cliClient.GetSecret(testSecretPath)

		assert.ErrorIs(t, expectedError, err)
	})

}
