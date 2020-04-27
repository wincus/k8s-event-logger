package main

import (
	"log"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	syncInterval time.Duration = 1 * time.Minute
)

func main() {

	log.Printf("starting k8s-event-logger")

	config, err := rest.InClusterConfig()

	if err != nil {
		log.Printf("could not instantiate kubernetes config: %v", err)
		os.Exit(1)
	}

	client, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Printf("could not instantiate a kubernetes client: %v", err)
		os.Exit(1)
	}

	factory := informers.NewSharedInformerFactory(client, syncInterval)
	podInformer := factory.Core().V1().Pods()
	endpointInformer := factory.Core().V1().Endpoints()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    logAdd,
		UpdateFunc: logUpdate,
		DeleteFunc: logDelete,
	})

	endpointInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    logAdd,
		UpdateFunc: logUpdate,
		DeleteFunc: logDelete,
	})

	factory.Start(wait.NeverStop)
	factory.WaitForCacheSync(wait.NeverStop)

	time.Sleep(10 * time.Minute)

	<-wait.NeverStop
}

func logAdd(new interface{}) {

	switch new := new.(type) {
	case *v1.Pod:
		log.Printf("NEW Pod Event: %v", new.Name)
	case *v1.Endpoints:
		log.Printf("NEW Endpoint Event: %v", new.Name)
	}

	return
}

func logUpdate(old, new interface{}) {

	switch new := new.(type) {
	case *v1.Pod:
		log.Printf("UPDATE Pod Event: %v", new.Name)
	case *v1.Endpoints:
		log.Printf("UPDATE Endpoint Event: %v", new.Name)
	}

	return
}

func logDelete(obj interface{}) {

	switch obj := obj.(type) {
	case *v1.Pod:
		log.Printf("DELETE Pod Event: %v", obj.Name)
	case *v1.Endpoints:
		log.Printf("DELETE Endpoint Event: %v", obj.Name)
	}

	return
}
