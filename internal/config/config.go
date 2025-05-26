package config

import (
	"bytes"
	"maps"
	"os"
	"sort"

	"github.com/goccy/go-yaml"
)

type RepositoryConfiguration map[string]string
type Configuration struct {
	rawConfig    map[string]RepositoryConfiguration
	Repositories []string
}

var readFileFunc = os.ReadFile

func (c Configuration) GetConfigurationForRepository(repository string) RepositoryConfiguration {
	merged := make(RepositoryConfiguration)

	maps.Copy(merged, c.rawConfig["common"])
	maps.Copy(merged, c.rawConfig[repository])

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
	if err = dec.Decode(&config.rawConfig); err != nil {
		return nil, err
	}

	config.Repositories = extractRepositoryNamesFromConfig(config.rawConfig)

	return config, nil
}

func NewConfigFromFile(path string) (config *Configuration, err error) {
	configFile, err := readFileFunc(path)
	if err != nil {
		return nil, err
	}

	return NewConfigFromReader(bytes.NewReader(configFile))
}
