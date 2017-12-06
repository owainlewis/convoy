package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ConvoyController struct {
	client    kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.Controller
	podLister lister_v1.PodLister
}

// NewConvoyController creates a new Convoy controller
func NewConvoyController(client kubernetes.Interface) *ConvoyController {

	namespace := "default"

	ctrl := &ConvoyController{
		client: client,
		queue:  workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	indexer, informer := cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(lo meta_v1.ListOptions) (runtime.Object, error) {
				return client.Core().Pods(namespace).List(lo)
			},
			WatchFunc: func(lo meta_v1.ListOptions) (watch.Interface, error) {
				return client.Core().Pods(namespace).Watch(lo)
			},
		},
		&v1.Pod{},
		10*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				glog.Infof("Adding item to queue: %s", obj)

				if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
					ctrl.queue.Add(key)
				}

			},
			UpdateFunc: func(old, new interface{}) {
				glog.Info("Updated object")
			},
			DeleteFunc: func(obj interface{}) {
				glog.Info("Deleted object")
			},
		},
		cache.Indexers{},
	)

	ctrl.informer = informer
	ctrl.podLister = lister_v1.NewPodLister(indexer)

	return ctrl
}

func (c *ConvoyController) Run(stopCh chan struct{}) {
	defer c.queue.ShutDown()
	glog.Info("Starting Controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		glog.Error(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
	glog.Info("Stopping Reboot Controller")
}

func (c *ConvoyController) runWorker() {
	for c.processNext() {
	}
}

func (c *ConvoyController) processNext() bool {
	// Wait until there is a new item in the working queue
	key, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(key)

	// Process the item. TODO handle errors in processing item
	c.process(key.(string))
	return true
}

func (c *ConvoyController) process(key string) error {
	glog.Info("Processing item from work queue")
	return nil
}
