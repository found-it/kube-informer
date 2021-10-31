package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	logger    *logrus.Entry
	clientset kubernetes.Interface
	factory   informers.SharedInformerFactory
	informer  cache.SharedInformer
	queue     workqueue.RateLimitingInterface
}

type Item struct {
	key       string
	action    string
	name      string
	namespace string
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

			c.logger.Infof("%s: %s", item.action, item.key)
			c.queue.Forget(obj)
			return nil
		}(obj)

		if err != nil {
			runtime.HandleError(err)
		}
	}
}

func (c *Controller) Run() {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	stopper := make(chan struct{})
	defer close(stopper)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	c.logger.Info("Listening for signals")

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
	s := <-sig
	c.logger.Info("Received signals ", s)
	c.logger.Info("Finished working")
}

func newController(clientset kubernetes.Interface, resyncPeriod time.Duration) *Controller {

	factory := informers.NewSharedInformerFactory(clientset, resyncPeriod)

	c := &Controller{
		logger:    logrus.WithField("app", "kai"),
		clientset: clientset,
		factory:   factory,
		informer:  factory.Core().V1().Pods().Informer(),
		queue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
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
			if err == nil {
				c.queue.Add(item)
			} else {
				logrus.Errorf("DELETE %s", pod.GetName())
			}
		},
	})
	return c
}

type arguments struct {
	kubeconfig string
}

func getArguments() arguments {
	var args arguments

	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&args.kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&args.kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	return args
}

// func doTheThing() {
// 	fmt.Println("Iter")
// }

func main() {
	logrus.Info("Shared Informer app started")

	arg := getArguments()
	config, err := clientcmd.BuildConfigFromFlags("", arg.kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	c := newController(clientset, 0)

	c.Run()

	// wait.Forever(doTheThing, 5*time.Second)
}
