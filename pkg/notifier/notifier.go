package notifier

import (
	v1 "k8s.io/api/core/v1"
)

type Notifier interface {
	Dispatch(event *v1.Event) error
}
