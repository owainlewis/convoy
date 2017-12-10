package dispatch

import (
	"fmt"

	glog "github.com/golang/glog"
	slack "github.com/nlopes/slack"
	v1 "k8s.io/api/core/v1"
)

// SlackNotifier is a notifier that will dispatch Kubernetes Events
// to a Slack channel
type SlackNotifier struct {
	Client  *slack.Client
	Channel string
}

// NewSlackNotifier will build a notifier that dispatches to a Slack channel
func NewSlackNotifier(token, channel string) *SlackNotifier {
	client := slack.New(token)
	return &SlackNotifier{
		Client:  client,
		Channel: channel,
	}
}

// Dispatch sends an event to Slack
func (n *SlackNotifier) Dispatch(event *v1.Event) error {
	msg := formatEventMessageForSlack(event)
	glog.Infof("Event: %s", msg)
	_, _, err := n.Client.PostMessage(n.Channel, msg, slack.PostMessageParameters{})
	return err
}

func formatEventMessageForSlack(event *v1.Event) string {
	return fmt.Sprintf("%s: %s", event.InvolvedObject.Kind, event.Message)
}
