package notifier

import (
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
)

// ConsoleNotifier is a notifier that writes to the console
type ConsoleNotifier struct {
}

// NewConsoleNotifier instantiates a default console notifier
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

// Dispatch will write a Kubernetes Event to the console
func (n *ConsoleNotifier) Dispatch(event *v1.Event) {
	glog.Infof("Event: %v", event)
}
