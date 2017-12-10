package config

import "testing"

const validConfig = `
slack:
  enabled: true
  channel: general
`

func TestUnmarshallValidConfig(t *testing.T) {
	conf, err := unmarshall([]byte(validConfig))
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !conf.Slack.Enabled {
		t.Errorf("Expected Slack to be enabled but was %v", conf.Slack.Enabled)
	}

	if conf.Slack.Channel != "general" {
		t.Errorf("Expected Slack channel to eq 'general' but got %v", conf.Slack.Channel)
	}
}
