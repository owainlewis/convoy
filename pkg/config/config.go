package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Filter defines the Kube event filters to be applied while watching events
type Filter struct {
	Type string
}

// Config defines the user level configuration options
type Config struct {
	Slack struct {
		Enabled bool
		Channel string
	}
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
