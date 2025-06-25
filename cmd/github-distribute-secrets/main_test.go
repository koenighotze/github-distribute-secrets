package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"koenighotze.de/github-distribute-secrets/internal/config"
	"koenighotze.de/github-distribute-secrets/pkg/github"
	"koenighotze.de/github-distribute-secrets/pkg/onepassword"
)

func TestMain(t *testing.T) {
	calledNewGhClientWithValue := false
	calledGithubSecretDistribution := false

	originalArgs := os.Args
	orignalMyNewGhClient := myNewGhClient
	originalMyGithubSecretDistribution := myGithubSecretDistribution

	defer func() {
		myNewGhClient = orignalMyNewGhClient
		myGithubSecretDistribution = originalMyGithubSecretDistribution
		os.Args = originalArgs
	}()

	os.Args = []string{"cmd"}
	myNewGhClient = func(dryRun bool) github.GithubClient {
		calledNewGhClientWithValue = dryRun
		return &mockGithubClient{}
	}
	myGithubSecretDistribution = func(configFileReader config.ConfigFileReader, op onepassword.OnePasswordClient, gh github.GithubClient, dumpConfig bool) bool {
		calledGithubSecretDistribution = true
		return true
	}

	t.Run("should use the dry run client if the flag is provided", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"cmd", "--dry-run"}

		main()

		assert.True(t, calledNewGhClientWithValue)
	})

	t.Run("should pass dump=true to githubSecretDistribution when --dump-config flag is provided", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"cmd", "--dump-config"}

		dumpFlagValue := false
		myGithubSecretDistribution = func(configFileReader config.ConfigFileReader, op onepassword.OnePasswordClient, gh github.GithubClient, dumpConfig bool) bool {
			dumpFlagValue = dumpConfig
			return true
		}

		main()

		assert.True(t, dumpFlagValue, "Should pass true for dump flag when --dump-config is provided")
	})

	t.Run("should pass dump=false to githubSecretDistribution when --dump-config flag is not provided", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"cmd"}

		dumpFlagValue := true
		myGithubSecretDistribution = func(configFileReader config.ConfigFileReader, op onepassword.OnePasswordClient, gh github.GithubClient, dumpConfig bool) bool {
			dumpFlagValue = dumpConfig
			return true
		}

		main()

		assert.False(t, dumpFlagValue, "Should pass false for dump flag when --dump-config is not provided")
	})

	t.Run("should use the default client if the flag is omitted", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = []string{"cmd"}

		main()

		assert.False(t, calledNewGhClientWithValue)

	})

	t.Run("should distribute the secrets", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()

		assert.True(t, calledGithubSecretDistribution)
	})
	t.Run("should exit with -1 if setting the secrets fails", func(t *testing.T) {
		t.Skip("skipping until we can test exit in a sane way")
	})
}
