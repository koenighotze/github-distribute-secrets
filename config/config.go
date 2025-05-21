package config

import (
	"bytes"
	"maps"
	"os"

	"github.com/goccy/go-yaml"
)

type RepositoryConfiguration map[string]string
type Configuration struct {
	rawConfig    map[string]RepositoryConfiguration
	Repositories []string
}

func (c Configuration) GetConfigurationForRepository(repository string) RepositoryConfiguration {
	merged := make(RepositoryConfiguration)

	maps.Copy(merged, c.rawConfig["common"])
	maps.Copy(merged, c.rawConfig[repository])

	return merged
}

func NewConfigFromReader(reader *bytes.Reader) (config *Configuration, err error) {
	config = &Configuration{}
	dec := yaml.NewDecoder(reader)
	if err = dec.Decode(&config.rawConfig); err != nil {
		return nil, err
	}

	for key := range maps.Keys(config.rawConfig) {
		if key == "common" {
			continue
		}

		config.Repositories = append(config.Repositories, key)
	}

	return config, nil
}

func NewConfigFromFile(path string) (config *Configuration, err error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewConfigFromReader(bytes.NewReader(configFile))
}
