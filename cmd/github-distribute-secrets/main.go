package main

import (
	"fmt"
	"log"

	"koenighotze.de/github-distribute-secrets/internal/config"
	"koenighotze.de/github-distribute-secrets/pkg/github"
	"koenighotze.de/github-distribute-secrets/pkg/onepassword"
)

func main() {
	githubSecretDistribution(config.NewConfigFileReader(), onepassword.NewClient(), github.NewClient())
}

func githubSecretDistribution(configFileReader config.ConfigFileReader, op onepassword.OnePasswordClient, gh github.GithubClient) {
	fmt.Println("Github Secret Distribution")

	configuration, err := configFileReader.ReadConfiguration("./config.yml")
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	if allOk := applyConfiguration(configuration, op, gh); !allOk {
		log.Default().Fatalln("Configuration was not applied successfully!")
	}
}

func applyConfigurationToRepository(configMap config.RepositoryConfiguration, repositoy string, op onepassword.OnePasswordClient, gh github.GithubClient) (ok bool) {
	ok = true

	for key, onePasswordPath := range configMap {
		secret, err := op.GetSecret(onePasswordPath)
		if err != nil {
			log.Default().Printf("Error reading secret %s: %v", key, err)
			ok = false
			continue
		}

		if err = gh.AddSecretToRepository(key, secret, repositoy); err != nil {
			log.Default().Printf("Error adding secret with key %s to repository %s: %v", key, repositoy, err)
			ok = false
		}
	}

	return ok
}

func applyConfiguration(configuration *config.Configuration, op onepassword.OnePasswordClient, gh github.GithubClient) (allOk bool) {
	allOk = true
	for _, repository := range configuration.Repositories {
		if ok := applyConfigurationToRepository(configuration.GetConfigurationForRepository(repository), repository, op, gh); !ok {
			log.Default().Printf("Cannot apply config to repository %s successfully!", repository)
			allOk = false
		}
	}
	return
}
