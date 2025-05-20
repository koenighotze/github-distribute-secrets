package main

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"os"
	"os/exec"
	"strings"

	"github.com/goccy/go-yaml"
)

type configMap map[string]string
type configuration map[string]configMap
type secretCacheType map[string]string

func main() {
	fmt.Println("Github Secret Distribution")

	config, err := readConfigurationFromFile()
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %w", err))
	}

	if allOk := applyConfiguration(config); !allOk {
		log.Default().Fatalln("Configuration was not applied successfully!")
	}
}

func getMergedConfig(config configuration, repository string) (merged configMap) {
	merged = make(configMap)

	maps.Copy(merged, config["common"])
	maps.Copy(merged, config[repository])

	return
}

func readSecretFromOnePassword(secretPath string, cache secretCacheType) (secret string, err error) {
	if cachedSecret, exists := cache[secretPath]; exists {
		log.Default().Printf("Using cached value for secret: %s", secretPath)
		return cachedSecret, nil
	}

	out, err := exec.Command("op", "read", secretPath).Output()
	if err != nil {
		log.Default().Printf("Error reading secret: %s", err)
		return
	}

	secret = strings.TrimSpace(string(out))
	cache[secretPath] = secret

	return
}

func applyConfigurationToRepository(configMap configMap, repositoy string, cache secretCacheType) (ok bool) {
	for key, onePasswordPath := range configMap {
		log.Default().Printf("Reading key %s with path %s", key, onePasswordPath)

		secret, err := readSecretFromOnePassword(onePasswordPath, cache)
		if err != nil {
			log.Default().Printf("Error reading secret %s: %v", key, err)
			continue
		}

		addSecretToRepository(key, secret, repositoy)
	}

	return true
}

func addSecretToRepository(key string, secret string, repositoy string) (err error) {
	log.Default().Printf("Adding secret %s to repository %s", key, repositoy)
	cmd := exec.Command("gh", "secret", "set", key, "--body", secret, "--repo", repositoy)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Default().Printf("Error adding secret %s to repository %s: %v\nOutput: %s",
			key, repositoy, err, string(output))
		return
	}

	log.Default().Printf("Successfully added secret %s to repository %s", key, repositoy)
	return nil
}

func applyConfiguration(config configuration) (allOk bool) {
	cache := make(secretCacheType)

	allOk = true
	for repositoy := range config {
		if repositoy == "common" {
			continue
		}

		log.Default().Printf("Handling repository: %s\n", repositoy)

		mergedConfig := getMergedConfig(config, repositoy)

		if ok := applyConfigurationToRepository(mergedConfig, repositoy, cache); !ok {
			log.Default().Printf("Cannot apply config to repository %s successfully!", repositoy)
			allOk = false
		}
	}
	return
}

func readConfigurationFromFile() (config configuration, err error) {
	configFile, err := os.ReadFile("./config.yml")
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(configFile))
	if err = dec.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
