package github

import (
	"log"

	"koenighotze.de/github-distribute-secrets/pkg/cli"
)

type Github interface {
	AddSecretToRepository(key string, secret string, repositoy string) (err error)
}

type GithubClient struct {
	runner cli.CommandRunner
}

func (gh *GithubClient) AddSecretToRepository(key string, secret string, repositoy string) (err error) {
	log.Default().Printf("In repository %s. Adding secret %s", repositoy, key)
	output, err := gh.runner.Run("gh", "secret", "set", key, "--body", secret, "--repo", repositoy)
	if err != nil {
		log.Default().Printf("Error adding secret %s to repository %s: %v\nOutput: %s",
			key, repositoy, err, string(output))
		return
	}
	return nil
}

func NewClient() Github {
	return &GithubClient{
		runner: cli.NewCommandRunner(),
	}
}
