package onepassword

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("should return the caching client", func(t *testing.T) {
		result := NewClient()

		_, ok := result.(cachedClient)

		assert.True(t, ok, "Expected result to be of type cachedClient")
	})
}

type mockRunner struct {
	Output []byte
	Err    error
}

func (m mockRunner) Run(name string, args ...string) ([]byte, error) {
	return m.Output, m.Err
}

func TestGetSecret(t *testing.T) {
	t.Run("should return the secret from onepassword", func(t *testing.T) {
		client := cliClient{
			runner: mockRunner{Output: []byte("supersecret\n")},
		}

		result, err := client.GetSecret("somepath")

		assert.Nil(t, err)
		assert.Equal(t, "supersecret", result)
	})

	t.Run("should trim whitespaces from the secret", func(t *testing.T) {
		client := cliClient{
			runner: mockRunner{Output: []byte("   supersecret   \n")},
		}

		result, err := client.GetSecret("somepath")

		assert.Nil(t, err)
		assert.Equal(t, "supersecret", result)
	})
	t.Run("should return the error if reading from onepassword fails", func(t *testing.T) {
		client := cliClient{
			runner: mockRunner{Output: nil, Err: errors.New("bumm")},
		}

		_, err := client.GetSecret("somepath")

		assert.ErrorContains(t, err, "bumm")
	})
}

func TestGetSecretWithCache(t *testing.T) {
	prepareClient := func(path string, output []byte, err error, cache *CacheEntry) cachedClient {
		cliClient := cachedClient{
			Cache: make(SecretCacheType),
			Op: cliClient{
				runner: mockRunner{Output: output, Err: err},
			},
		}
		if nil != cache {
			cliClient.Cache[path] = *cache
		}

		return cliClient
	}

	t.Run("should return the uncached value if uncached", func(t *testing.T) {
		cliClient := prepareClient("somepath", []byte("UncachedOutput"), nil, nil)

		result, _ := cliClient.GetSecret("somepath")

		assert.Equal(t, "UncachedOutput", result)

	})

	t.Run("should return the cached value if cached", func(t *testing.T) {
		cliClient := prepareClient("somepath", []byte("UncachedOutput"), nil, &CacheEntry{Value: "cached", Err: nil})

		result, _ := cliClient.GetSecret("somepath")

		assert.Equal(t, "cached", result)
	})

	t.Run("should return the uncached error if uncached", func(t *testing.T) {
		cliClient := prepareClient("somepath", nil, assert.AnError, nil)

		_, err := cliClient.GetSecret("somepath")

		assert.ErrorIs(t, err, assert.AnError)
	})
	t.Run("should return the cached error if an error occured", func(t *testing.T) {
		expectedError := errors.New("cachederror")
		cliClient := prepareClient("somepath", nil, assert.AnError, &CacheEntry{Value: "cached", Err: expectedError})

		_, err := cliClient.GetSecret("somepath")

		assert.ErrorIs(t, expectedError, err)
	})

}
