package main

import (
	"flag"
	"log"

	"koenighotze.de/github-distribute-secrets/internal/config"
	"koenighotze.de/github-distribute-secrets/pkg/github"
	"koenighotze.de/github-distribute-secrets/pkg/onepassword"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate execution without making changes")
	flag.Parse()
	if *dryRun {
		log.Println("RUNNING IN DRY-RUN MODE - Will not change anything!")
	}

	gh := github.NewClient(*dryRun)
	op := onepassword.NewClient()

	if !githubSecretDistribution(config.NewConfigFileReader(), op, gh) {
		log.Default().Fatalln("Not all configuration was applied successfully!")
	}
}
