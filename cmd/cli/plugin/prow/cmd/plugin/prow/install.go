package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"

	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install Prow, its repo, its packages, and prerequisite secrets",
	Args:  cobra.NoArgs,
	Example: `
	tanzu prow install`,
	RunE: installProw,
}

func installProw(cmd *cobra.Command, _ []string) error {
	var (
		kubeConfig = getDefaultKubeconfigPath()
	)

	clientset := getClientSet(kubeConfig)
	fmt.Println(clientset)

	fmt.Println("One day, some day soon? Install Prow, the repo, all the repo package bundles, and its prerequisites on a workload cluster.")
	return nil
}

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
