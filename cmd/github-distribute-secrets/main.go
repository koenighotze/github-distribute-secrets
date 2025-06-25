package main

import (
	"flag"
	"log"

	"koenighotze.de/github-distribute-secrets/internal/config"
	"koenighotze.de/github-distribute-secrets/pkg/github"
	"koenighotze.de/github-distribute-secrets/pkg/onepassword"
)

var (
	myNewGhClient              = github.NewClient
	myNewOpClient              = onepassword.NewClient
	myNewConfigFileReader      = config.NewConfigFileReader
	myGithubSecretDistribution = githubSecretDistribution
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate execution without making changes")
	dumpConfig := flag.Bool("dump-config", false, "Dump configuration without applying it")
	flag.Parse()

	if *dryRun {
		log.Println("RUNNING IN DRY-RUN MODE - Will not change anything!")
	}

	if *dumpConfig {
		log.Println("CONFIGURATION DUMP ENABLED - Configuration will be printed")
	}

	gh := myNewGhClient(*dryRun)
	op := myNewOpClient()

	if !myGithubSecretDistribution(myNewConfigFileReader(), op, gh, *dumpConfig) {
		log.Default().Fatalln("Not all configuration was applied successfully!")
	}
}
