package onepassword

import (
	"log"
	"strings"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type cacheEntry struct {
	Value string
	Err   error
}

type secretCacheType map[string]cacheEntry

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
	Cache secretCacheType
	Op    OnePasswordClient
}

func (c *cachedClient) GetSecret(secretPath string) (secret string, err error) {
	if cachedSecret, exists := c.Cache[secretPath]; exists {
		return cachedSecret.Value, cachedSecret.Err
	}

	secret, err = c.Op.GetSecret(secretPath)
	if err == nil {
		c.Cache[secretPath] = cacheEntry{secret, nil}
	} else {
		c.Cache[secretPath] = cacheEntry{"", err}
	}

	return
}

func NewClient() OnePasswordClient {
	client := &cachedClient{
		Cache: make(secretCacheType),
		Op: &cliClient{
			runner: cli.NewCommandRunner(),
		},
	}

	return client
}
