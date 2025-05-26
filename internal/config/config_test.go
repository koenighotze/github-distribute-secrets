package config

import (
	"bytes"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	yamlConfigurationFull = `
common:
   KEY0: VAL0
repo1:
   KEY1: VAL1
repo2:
   KEY2: VAL2
`
	yamlConfigurationOverwrites = `
common:
   KEY0: VAL0
repo1:
   KEY1: VAL1
repo2:
   KEY2: VAL2
`

	yamlConfigurationCommonOnly = `
common:
   KEY0: VAL0
`

	yamlConfigurationNoCommon = `
repo1:
   KEY1: VAL1
`
)

func TestExtractingRepositoriesFromConfiguration(t *testing.T) {
	t.Run("should return an empty array if configration is empty", func(t *testing.T) {
		rawConfig := make(map[string]RepositoryConfiguration)

		result := extractRepositoryNamesFromConfig(rawConfig)

		assert.Equal(t, 0, len(result))
	})

	t.Run("should return an empty array if only the common section is found", func(t *testing.T) {
		rawConfig := map[string]RepositoryConfiguration{
			"common": {"foo": "bar"},
		}

		result := extractRepositoryNamesFromConfig(rawConfig)

		assert.Equal(t, 0, len(result))
	})

	t.Run("should return the repositories", func(t *testing.T) {
		rawConfig := map[string]RepositoryConfiguration{
			"common": {"foo": "bar"},
			"bar":    {"foo": "bar"},
			"qux":    {"foo": "bar"},
		}

		result := extractRepositoryNamesFromConfig(rawConfig)

		assert.Equal(t, []string{"bar", "qux"}, result)
	})
}

func TestNewConfigFromReader(t *testing.T) {
	t.Run("should return the common configration", func(t *testing.T) {
		expectedConfig := &Configuration{
			Repositories: []string{},
			rawConfig: map[string]RepositoryConfiguration{
				"common": map[string]string{"KEY0": "VAL0"},
			},
		}

		reader := bytes.NewReader([]byte(yamlConfigurationCommonOnly))
		result, err := NewConfigFromReader(reader)

		assert.Nil(t, err)
		assert.Equal(t, expectedConfig, result)
	})

	t.Run("should initialize the repositories", func(t *testing.T) {
		expectedConfig := &Configuration{
			Repositories: []string{"repo1", "repo2"},
			rawConfig: map[string]RepositoryConfiguration{
				"common": map[string]string{"KEY0": "VAL0"},
				"repo1":  map[string]string{"KEY1": "VAL1"},
				"repo2":  map[string]string{"KEY2": "VAL2"},
			},
		}

		reader := bytes.NewReader([]byte(yamlConfigurationFull))
		result, err := NewConfigFromReader(reader)
		sort.Strings(result.Repositories)

		assert.Nil(t, err)
		assert.Equal(t, expectedConfig, result)
	})

	t.Run("should return the error if the configuration was invalid", func(t *testing.T) {
		reader := bytes.NewReader([]byte("something invalid"))
		result, err := NewConfigFromReader(reader)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGetConfigurationForRepository(t *testing.T) {
	t.Run("should return the configuration containing common and repo", func(t *testing.T) {
		reader := bytes.NewReader([]byte(yamlConfigurationFull))
		config, _ := NewConfigFromReader(reader)

		result := config.GetConfigurationForRepository("repo1")

		assert.Equal(t, RepositoryConfiguration{
			"KEY0": "VAL0",
			"KEY1": "VAL1",
		}, result)
	})

	t.Run("should return the configuration without the fields of other repos", func(t *testing.T) {
		reader := bytes.NewReader([]byte(yamlConfigurationFull))
		config, _ := NewConfigFromReader(reader)

		result := config.GetConfigurationForRepository("repo1")

		assert.NotContains(t, result, "KEY2")
	})

	t.Run("should return the common section if the repo is not found", func(t *testing.T) {
		reader := bytes.NewReader([]byte(yamlConfigurationFull))
		config, _ := NewConfigFromReader(reader)

		result := config.GetConfigurationForRepository("not there")

		assert.Equal(t, RepositoryConfiguration{
			"KEY0": "VAL0",
		}, result)
	})

	t.Run("should return an repo configuration if nothing common is defined", func(t *testing.T) {
		reader := bytes.NewReader([]byte(yamlConfigurationNoCommon))
		config, _ := NewConfigFromReader(reader)

		result := config.GetConfigurationForRepository("repo1")

		assert.Equal(t, RepositoryConfiguration{
			"KEY1": "VAL1",
		}, result)
	})

	t.Run("should override common values with repo level settings", func(t *testing.T) {
		reader := bytes.NewReader([]byte(`
common:
   A: B
   R: T
repo1:
   A: C
   B: D
`))
		config, _ := NewConfigFromReader(reader)

		result := config.GetConfigurationForRepository("repo1")

		assert.Equal(t, RepositoryConfiguration{
			"A": "C",
			"B": "D",
			"R": "T",
		}, result)
	})
}

func TestNewConfigFromFile(t *testing.T) {
	t.Run("should return the error if reading the file fails", func(t *testing.T) {
		orgFileFunc := readFileFunc
		readFileFunc = func(name string) ([]byte, error) {
			return nil, os.ErrExist
		}
		defer func() { readFileFunc = orgFileFunc }()

		_, err := NewConfigFromFile("somefile")

		assert.Equal(t, os.ErrExist, err)
	})
	t.Run("should return the configuration", func(t *testing.T) {
		orgFileFunc := readFileFunc
		readFileFunc = func(name string) ([]byte, error) {
			return []byte(yamlConfigurationFull), nil
		}
		defer func() { readFileFunc = orgFileFunc }()

		expectedConfig := &Configuration{
			Repositories: []string{"repo1", "repo2"},
			rawConfig: map[string]RepositoryConfiguration{
				"common": map[string]string{"KEY0": "VAL0"},
				"repo1":  map[string]string{"KEY1": "VAL1"},
				"repo2":  map[string]string{"KEY2": "VAL2"},
			},
		}

		result, err := NewConfigFromFile("somefile")
		sort.Strings(result.Repositories)

		assert.Nil(t, err)
		assert.Equal(t, expectedConfig, result)

	})

	t.Run("should return the error if the file contains non-yaml data", func(t *testing.T) {
		orgFileFunc := readFileFunc
		readFileFunc = func(name string) ([]byte, error) {
			return []byte("ffff"), nil
		}
		defer func() { readFileFunc = orgFileFunc }()

		_, err := NewConfigFromFile("somefile")

		assert.Contains(t, err.Error(), "mapping")
	})

}
