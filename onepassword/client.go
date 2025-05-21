package onepassword

import (
	"log"
	"os/exec"
	"strings"
)

type CacheEntry struct {
	Value string
	Err   error
}

type SecretCacheType map[string]CacheEntry

type OnePasswordClient interface {
	GetSecret(secretPath string) (secret string, err error)
}

type cliClient struct{}

func (d cliClient) GetSecret(secretPath string) (secret string, err error) {
	log.Default().Printf("Reading secret with path %s", secretPath)

	out, err := exec.Command("op", "read", secretPath).Output()
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

func (c cachedClient) GetSecret(secretPath string) (secret string, err error) {
	if cachedSecret, exists := c.Cache[secretPath]; exists {
		log.Default().Printf("Using cached value for secret: %s", secretPath)
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
	return cachedClient{
		Cache: make(SecretCacheType),
		Op:    cliClient{},
	}
}
