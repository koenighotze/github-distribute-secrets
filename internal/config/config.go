package config

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"sort"

	"github.com/goccy/go-yaml"
)

type RepositoryConfiguration map[string]string
type Configuration struct {
	RawConfig    map[string]RepositoryConfiguration
	Repositories []string
}
type ConfigFileReader interface {
	ReadConfiguration(path string) (config *Configuration, err error)
}

func (c Configuration) GetConfigurationForRepository(repository string) RepositoryConfiguration {
	merged := make(RepositoryConfiguration)

	maps.Copy(merged, c.RawConfig["common"])
	maps.Copy(merged, c.RawConfig[repository])

	return merged
}

func extractRepositoryNamesFromConfig(rawConfig map[string]RepositoryConfiguration) []string {
	result := make([]string, 0, len(rawConfig))
	for key := range maps.Keys(rawConfig) {
		if key == "common" {
			continue
		}

		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func NewConfigFromReader(reader *bytes.Reader) (config *Configuration, err error) {
	config = &Configuration{}
	dec := yaml.NewDecoder(reader)
	if err = dec.Decode(&config.RawConfig); err != nil {
		return nil, err
	}

	config.Repositories = extractRepositoryNamesFromConfig(config.RawConfig)

	return config, nil
}

type configFileReader struct {
	fileReader func(name string) ([]byte, error)
}

func (reader configFileReader) ReadConfiguration(path string) (config *Configuration, err error) {
	configFile, err := reader.fileReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	return NewConfigFromReader(bytes.NewReader(configFile))
}

func NewConfigFileReader() ConfigFileReader {
	return configFileReader{
		fileReader: os.ReadFile,
	}
}

// DumpConfiguration returns a formatted string representation of the configuration
// that will be applied, showing which secrets will be set for each repository.
func (c Configuration) DumpConfiguration() string {
	var buffer bytes.Buffer

	buffer.WriteString("Configuration Summary:\n")
	buffer.WriteString("=====================\n\n")

	// First show common secrets
	commonConfig := c.RawConfig["common"]
	if len(commonConfig) > 0 {
		buffer.WriteString("Common Secrets (applied to all repositories):\n")
		for key := range commonConfig {
			buffer.WriteString(fmt.Sprintf("  - %s\n", key))
		}
		buffer.WriteString("\n")
	}

	// Then show per-repository secrets
	buffer.WriteString("Repository-Specific Configurations:\n")
	for _, repo := range c.Repositories {
		repoConfig := c.GetConfigurationForRepository(repo)
		buffer.WriteString(fmt.Sprintf("- %s:\n", repo))

		if len(repoConfig) == 0 {
			buffer.WriteString("  No secrets configured\n")
		} else {
			for key := range repoConfig {
				buffer.WriteString(fmt.Sprintf("  - %s\n", key))
			}
		}
	}

	return buffer.String()
}
