package controller

import (
	"fmt"
	"log"
	"time"

	glog "github.com/golang/glog"
	notifier "github.com/owainlewis/convoy/pkg/notifier"
	v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/util/runtime"
	wait "k8s.io/apimachinery/pkg/util/wait"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	workqueue "k8s.io/client-go/util/workqueue"
)

const (
	// ConvoyEventType defines the type of event to watch
	ConvoyEventType = "Pod"
)

// ConvoyController defines the structure of the controller
type ConvoyController struct {
	client            kubernetes.Interface
	eventGetter       corev1.EventsGetter
	eventLister       listerv1.EventLister
	eventListerSynced cache.InformerSynced
	queue             workqueue.RateLimitingInterface
	notifier          notifier.Notifier
}

// NewConvoyController creates a new Convoy controller
func NewConvoyController(client kubernetes.Interface, informer informercorev1.EventInformer, notifier notifier.Notifier) *ConvoyController {
	c := &ConvoyController{
		client:            client,
		eventGetter:       client.CoreV1(),
		eventLister:       informer.Lister(),
		eventListerSynced: informer.Informer().HasSynced,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		notifier:          notifier,
	}

	informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				c.enqueue(obj)
			},
		},
	)

	return c
}

// Run will start the controller
func (c *ConvoyController) Run(stopCh chan struct{}) {
	defer c.queue.ShutDown()

	glog.Info("Waiting for cache sync")
	if !cache.WaitForCacheSync(stopCh, c.eventListerSynced) {
		glog.Info("Timeout waiting for caches to sync")
		return
	}
	log.Print("Caches are synced")

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
	glog.Info("Stopping controller")
}

func (c *ConvoyController) enqueue(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.queue.AddRateLimited(key)
}

func (c *ConvoyController) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *ConvoyController) processNextWorkItem() bool {
	obj, shutdown := c.queue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.queue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.queue.Forget(obj)

		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *ConvoyController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	event, err := c.eventLister.Events(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("foo '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	c.processEvent(event)

	return nil
}

func (c *ConvoyController) processEvent(event *v1.Event) {
	if event.InvolvedObject.Kind == ConvoyEventType {
		c.notifier.Dispatch(event)
		//now := meta_v1.Now()
		//if !event.FirstTimestamp.Before(&now) {
		// glog.Infof("Processing %s event %s", event.InvolvedObject.Kind, event.Message)
		//}
	}
}
