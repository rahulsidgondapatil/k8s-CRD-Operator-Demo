package main

import (
	"os"

	"github.com/google/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	clientset "github.com/rahulsidgondapatil/sample-customController/pkg/client/clientset/versioned"
)

func getKubeClients() (*kubernetes.Clientset, *clientset.Clientset) {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		logger.Fatalf("getClusterConfig: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	return kubeClient, exampleClient
}
