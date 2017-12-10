package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Filter defines the Kube event filters to be applied while watching events
type Filter struct {
	Type string
}

// SlackConfig defines the configuration needed to auth with slack
type SlackConfig struct {
	// Should be moved into a secret
	Token   string
	Channel string
}

// Config defines the user level configuration options
type Config struct {
	Notifier struct {
		Slack SlackConfig
	}
	Filters []Filter
}

// FromFile will load controller configuration from a file path
func FromFile(filepath string) (*Config, error) {
	contents, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	return unmarshall(contents)
}

func unmarshall(content []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
