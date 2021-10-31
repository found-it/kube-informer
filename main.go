package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/found-it/kube-informer/inform/controller"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

	c := controller.NewController(clientset, 0)

	c.Run()

	// wait.Forever(doTheThing, 5*time.Second)
}
