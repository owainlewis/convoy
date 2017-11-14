package controller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ConvoyController struct {
	client   kubernetes.Interface
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
}

// NewConvoyController instantiates a new Convoy controller
func NewConvoyController(client kubernetes.Interface) *ConvoyController {
	ctrl := &ConvoyController{
		client: client,
		queue:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	return ctrl
}
