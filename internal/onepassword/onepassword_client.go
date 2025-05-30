package onepassword

import (
	"log"
	"strings"

	"koenighotze.de/github-distribute-secrets/internal/common/cli"
)

type CacheEntry struct {
	Value string
	Err   error
}

type SecretCacheType map[string]CacheEntry

type OnePasswordClient interface {
	GetSecret(secretPath string) (secret string, err error)
}

type cliClient struct {
	runner cli.CommandRunner
}

func (d *cliClient) GetSecret(secretPath string) (secret string, err error) {
	out, err := d.runner.Run("op", "read", secretPath)
	if err != nil {
		log.Default().Printf("Error reading secret: %s", err)
		return
	}

	secret = strings.TrimSpace(string(out))

	return
}

type cachedClient struct {
	Cache SecretCacheType
	Op    OnePasswordClient
}

func (c *cachedClient) GetSecret(secretPath string) (secret string, err error) {
	if cachedSecret, exists := c.Cache[secretPath]; exists {
		return cachedSecret.Value, cachedSecret.Err
	}

	secret, err = c.Op.GetSecret(secretPath)
	if err == nil {
		c.Cache[secretPath] = CacheEntry{secret, nil}
	} else {
		c.Cache[secretPath] = CacheEntry{"", err}
	}

	return
}

func NewClient() OnePasswordClient {
	client := &cachedClient{
		Cache: make(SecretCacheType),
		Op: &cliClient{
			runner: cli.NewCommandRunner(),
		},
	}

	return client
}
