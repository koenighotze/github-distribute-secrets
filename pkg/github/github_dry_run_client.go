package github

import (
	"fmt"
	"log"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type dryRunGithubClient struct {
	runner cli.CommandRunner
}

func (gh *dryRunGithubClient) AddSecretToRepository(key string, secret string, repositoy string) (err error) {
	log.Default().Printf("DRY RUN: In repository %s. Should add secret with key %s", repositoy, key)
	if _, err = gh.runner.Run("gh", "repo", "view", repositoy); err != nil {
		return fmt.Errorf("repository %s does not seem to exist. %w", repositoy, err)
	}
	return nil
}

func withDryRun() GithubClient {
	return &dryRunGithubClient{
		runner: cli.NewCommandRunner(),
	}
}
