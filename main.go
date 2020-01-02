package main

import (
	"flag"
	"fmt"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/klog"
	"k8s.io/sample-controller/pkg/signals"

	"github.com/rahulsidgondapatil/sample-customController/conf"
	informers "github.com/rahulsidgondapatil/sample-customController/pkg/client/informers/externalversions"
)

var (
	masterURL  string
	kubeconfig string
)

// main code path
func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := conf.New()
	if err != nil {
		klog.Fatal(fmt.Sprintf("Error in loading configuration. Error:%s", err.Error()))
	}

	kubeClient, exampleClient := getKubeClients()

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(cfg, kubeClient, exampleClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		exampleInformerFactory.Customcontroller().V1alpha1().DepSvcResources())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
