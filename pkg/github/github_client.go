package github

import (
	"fmt"
	"log"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type GithubClient interface {
	AddSecretToRepository(key string, secret string, repository string) (err error)
}

type cliGithubClient struct {
	runner cli.CommandRunner
}

func (gh *cliGithubClient) AddSecretToRepository(key string, secret string, repository string) (err error) {
	log.Printf("In repository %s. Adding secret with key %s", repository, key)
	if _, err = gh.runner.Run("gh", "secret", "set", key, "--body", secret, "--repo", repository); err != nil {
		return fmt.Errorf("failed adding secret as key %s to repository %s: %w", key, repository, err)
	}
	return nil
}

func NewClient(dryRun bool) GithubClient {
	if dryRun {
		return withDryRun()
	}

	return &cliGithubClient{
		runner: cli.NewCommandRunner(),
	}
}
