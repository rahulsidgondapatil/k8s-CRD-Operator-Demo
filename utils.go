package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	clientset "github.com/rahulsidgondapatil/sample-customController/pkg/client/clientset/versioned"
)

func getKubeClients() (*kubernetes.Clientset, *clientset.Clientset) {
	// construct the path to resolve to `~/.kube/config`
	configPath := os.Getenv("KUBECONFIG")
	if configPath == "" {
		configPath = os.Getenv("HOME") + "/.kube/config"
	}

	var conf *rest.Config

	// In cluster config
	conf, err := rest.InClusterConfig()
	if err != nil {
		klog.Infof("err=", err)
	}

	if nil == conf {
		//Use Outside of cluster config
		conf, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			klog.Fatal(fmt.Sprintf("error in getting Kubeconfig: %v", err))
		}
	}

	kubeClient, err := kubernetes.NewForConfig(conf)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(conf)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	return kubeClient, exampleClient
}
