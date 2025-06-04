package github

import (
	"fmt"
	"log"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type GithubClient interface {
	AddSecretToRepository(key string, secret string, repositoy string) (err error)
}

type cliGithubClient struct {
	runner cli.CommandRunner
}

func (gh *cliGithubClient) AddSecretToRepository(key string, secret string, repositoy string) (err error) {
	log.Default().Printf("In repository %s. Adding secret %s", repositoy, key)
	if _, err = gh.runner.Run("gh", "secret", "set", key, "--body", secret, "--repo", repositoy); err != nil {
		return fmt.Errorf("failed adding secret as key %s to repository %s: %w", key, repositoy, err)
	}
	return nil
}

func NewClient() GithubClient {
	return &cliGithubClient{
		runner: cli.NewCommandRunner(),
	}
}
