package notifier

import (
	slack "github.com/nlopes/slack"
	v1 "k8s.io/api/core/v1"
)

// SlackNotifier is a notifier that will dispatch Kubernetes Events
// to a Slack channel
type SlackNotifier struct {
	Client  *slack.Client
	Channel string
}

func NewSlackNotifier(client *slack.Client, channel string) *SlackNotifier {
	return &SlackNotifier{
		Client:  client,
		Channel: channel,
	}
}

// Dispatch will send an event to Slack
func (n *SlackNotifier) Dispatch(event *v1.Event) error {
	params := slack.PostMessageParameters{}
	_, _, err := n.Client.PostMessage("general", event.Message, params)
	return err
}
