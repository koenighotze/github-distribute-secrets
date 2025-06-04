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
	flag.Parse()
	if *dryRun {
		log.Println("RUNNING IN DRY-RUN MODE - Will not change anything!")
	}

	gh := myNewGhClient(*dryRun)
	op := myNewOpClient()

	if !myGithubSecretDistribution(myNewConfigFileReader(), op, gh) {
		log.Default().Fatalln("Not all configuration was applied successfully!")
	}
}
