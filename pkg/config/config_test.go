package config

import "testing"

const validConfig = `
notifier:
  slack:
    token: XXX
    channel: general
filters:
  - objectKind: Pod
`

func TestUnmarshallValidConfig(t *testing.T) {
	conf, err := unmarshall([]byte(validConfig))
	if err != nil {
		t.Errorf("Expected valid config to be unmarshalled. Got %s", err)
	}

	slack := conf.Notifier.Slack
	if slack.Token != "XXX" {
		t.Errorf("Expected Slack token to be %s but got %s", "XXX", slack.Token)
	}
	if slack.Channel != "general" {
		t.Errorf("Expected Slack channel to be %s but got %s", "general", slack.Token)
	}
}
