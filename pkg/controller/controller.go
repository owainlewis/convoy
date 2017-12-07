package controller

import (
	"log"
	"time"

	glog "github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	wait "k8s.io/apimachinery/pkg/util/wait"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

type ConvoyController struct {
	client            kubernetes.Interface
	eventGetter       corev1.EventsGetter
	eventLister       listerv1.EventLister
	eventListerSynced cache.InformerSynced
	queue             workqueue.RateLimitingInterface
}

// NewConvoyController creates a new Convoy controller
func NewConvoyController(client kubernetes.Interface, informer informercorev1.EventInformer) *ConvoyController {
	c := &ConvoyController{
		client:            client,
		eventGetter:       client.CoreV1(),
		eventLister:       informer.Lister(),
		eventListerSynced: informer.Informer().HasSynced,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}

	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				event := obj.(*v1.Event)
				glog.Infof("New event %v", event.Message)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				event := newObj.(*v1.Event)
				glog.Infof("Update event %v", event.Message)
			},
		},
	)

	return c
}

func (c *ConvoyController) Run(stopCh chan struct{}) {

	defer c.queue.ShutDown()

	glog.Info("Starting controller")

	glog.Info("Waiting for cache sync")
	if !cache.WaitForCacheSync(stopCh, c.eventListerSynced) {
		glog.Info("Timeout waiting for caches to sync")
		return
	}
	log.Print("caches are synced")

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
	glog.Info("Stopping controller")
}

func (c *ConvoyController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *ConvoyController) processNextWorkItem() bool {
	return true
}
