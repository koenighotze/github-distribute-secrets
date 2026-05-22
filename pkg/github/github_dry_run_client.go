package github

import (
	"fmt"
	"log"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type dryRunGithubClient struct {
	runner cli.CommandRunner
}

func (gh *dryRunGithubClient) AddSecretToRepository(key string, secret string, repository string) (err error) {
	log.Printf("DRY RUN: In repository %s. Should add secret with key %s", repository, key)
	if _, err = gh.runner.Run("gh", "repo", "view", repository); err != nil {
		return fmt.Errorf("repository %s does not seem to exist. %w", repository, err)
	}
	return nil
}

func withDryRun() GithubClient {
	return &dryRunGithubClient{
		runner: cli.NewCommandRunner(),
	}
}
