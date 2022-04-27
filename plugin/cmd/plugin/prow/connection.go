package main

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getClientSet(kubeConfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func getDefaultKubeconfigPath() string {
	kubeConfigFilename := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	// fallback to default kubeconfig file location if no env variable set
	if kubeConfigFilename == "" {
		kubeConfigFilename = clientcmd.RecommendedHomeFile
	}
	return kubeConfigFilename
}
