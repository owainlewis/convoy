package notifier

import v1 "k8s.io/api/core/v1"

// SlackNotifier is a notifier that will dispatch Kubernetes Events
// to a Slack channel
type SlackNotifier struct {
	Channel string
}

// Dispatch will send an event to Slack
func (n *SlackNotifier) Dispatch(event *v1.Event) {

}
