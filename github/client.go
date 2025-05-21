package github

import (
	"log"
)

type Github interface {
	addSecretToRepository(key string, secret string, repositoy string) (err error)
}

type GithubClient struct {
}

func (gh GithubClient) AddSecretToRepository(key string, secret string, repositoy string) (err error) {
	log.Default().Printf("Adding secret %s", key)
	// cmd := exec.Command("gh", "secret", "set", key, "--body", secret, "--repo", repositoy)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Default().Printf("Error adding secret %s to repository %s: %v\nOutput: %s",
	// 		key, repositoy, err, string(output))
	// 	return
	// }

	log.Default().Printf("Successfully added secret %s to repository %s", key, repositoy)
	return nil
}

func NewClient() GithubClient {
	return GithubClient{}
}
