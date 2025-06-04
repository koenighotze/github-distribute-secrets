package github

import (
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
	output, err := gh.runner.Run("gh", "secret", "set", key, "--body", secret, "--repo", repositoy)
	if err != nil {
		log.Default().Printf("Error adding secret %s to repository %s: %v\nOutput: %s",
			key, repositoy, err, string(output))
		return
	}
	return nil
}

func NewClient() GithubClient {
	return &cliGithubClient{
		runner: cli.NewCommandRunner(),
	}
}
