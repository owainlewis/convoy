package dispatch

import (
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
)

// ConsoleNotifier is a dispatcher that writes to the console
type ConsoleNotifier struct {
}

// NewConsoleNotifier instantiates a default console dispatcher
func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

// Dispatch will write a Kubernetes Event to the console
func (n *ConsoleNotifier) Dispatch(event *v1.Event) error {
	glog.Infof("%s: %v", event.InvolvedObject.Kind, event.Message)
	return nil
}
