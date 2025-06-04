package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("should use the dry run client if the flag is provided", func(t *testing.T) {})
	t.Run("should use the default client if the flag is omitted", func(t *testing.T) {})
	t.Run("should distribute the secrets", func(t *testing.T) {})
	t.Run("should exit with -1 if setting the secrets fails", func(t *testing.T) {})
}
