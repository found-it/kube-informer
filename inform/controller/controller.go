package controller

import (
	"fmt"
	// "os"
	// "os/signal"
	// "syscall"
	"time"

	"github.com/found-it/kube-informer/inform/report"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	logger    *logrus.Entry
	clientset kubernetes.Interface
	factory   informers.SharedInformerFactory
	informer  cache.SharedInformer
	queue     workqueue.RateLimitingInterface
	report    *report.ReportController
}

type Item struct {
	key       string
	action    string
	name      string
	namespace string
	pod       *v1.Pod
}

func (c *Controller) worker() {
	for {
		obj, shutdown := c.queue.Get()

		if shutdown {
			c.logger.Info("Cache has been shut down")
			break
		}

		err := func(obj interface{}) error {
			defer c.queue.Done(obj)
			var item Item
			var ok bool

			if item, ok = obj.(Item); !ok {
				c.queue.Forget(obj)
				runtime.HandleError(fmt.Errorf("expected Item in queue but go %#v", obj))
				return nil
			}

			// c.logger.WithField("action", item.action).Infof("%s", item.key)
			switch item.action {
			case "ADD":
				c.report.Add(item.pod)
			case "UPDATE":
				c.report.Update(item.pod)
			case "DELETE":
				c.report.Delete(item.pod)
			}
			c.queue.Forget(obj)
			return nil
		}(obj)

		if err != nil {
			runtime.HandleError(err)
		}
	}
}

func (c *Controller) RunOnce() {
	runner(c, true)
}

func (c *Controller) Run() {
	runner(c, false)
}

func (c* Controller) printer() {
    for item := range c.report.Report.Items {
        c.logger.Info(item)
    }
}

func runner(c *Controller, runOnce bool) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	stopper := make(chan struct{})
	defer close(stopper)

	c.factory.Start(stopper)

	if !cache.WaitForCacheSync(stopper, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for cache to sync"))
		return
	}

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		c.logger.Infof("Launching worker%d", i)
		go wait.Until(c.worker, time.Second, stopper)
	}

	c.logger.Info("Started worker")
	if runOnce {
		time.Sleep(5 * time.Second)
		// build report from the cache and exit

	} else {
		// sig := make(chan os.Signal, 1)
		// signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
        //
		// c.logger.Info("Listening for signals")
		// s := <-sig
		// c.logger.Info("Received signals ", s)

        wait.Until(c.printer, 5*time.Second, stopper)
	}
	c.logger.Info("Finished working")
}

func NewController(clientset kubernetes.Interface, resyncPeriod time.Duration) *Controller {

	factory := informers.NewSharedInformerFactory(clientset, resyncPeriod)

	c := &Controller{
		logger:    logrus.WithField("pkg", "controller"),
		clientset: clientset,
		factory:   factory,
		informer:  factory.Core().V1().Pods().Informer(),
		queue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		report:    report.NewReportController(),
	}

	var item Item
	var err error

	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			item.key, err = cache.MetaNamespaceKeyFunc(obj)
			item.action = "ADD"
			item.name = pod.GetName()
			item.namespace = pod.GetNamespace()
			item.pod = pod
			if err == nil {
				c.queue.Add(item)
			} else {
				logrus.Errorf("ADD %s", pod.GetName())
			}
		},
		UpdateFunc: func(o, n interface{}) {
			oldpod := o.(*v1.Pod)
			newpod := n.(*v1.Pod)
			item.key, err = cache.MetaNamespaceKeyFunc(n)
			item.action = "UPDATE"
			item.name = newpod.GetName()
			item.namespace = newpod.GetNamespace()
			item.pod = newpod
			if err == nil {
				if oldpod.ResourceVersion != newpod.ResourceVersion {
					c.queue.Add(item)
				}
			} else {
				logrus.Errorf("UPDATE %s - %s", newpod.GetName(), oldpod.GetName())
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			item.key, err = cache.MetaNamespaceKeyFunc(obj)
			item.action = "DELETE"
			item.name = pod.GetName()
			item.namespace = pod.GetNamespace()
			item.pod = pod
			if err == nil {
				c.queue.Add(item)
			} else {
				logrus.Errorf("DELETE %s", pod.GetName())
			}
		},
	})
	return c
}
